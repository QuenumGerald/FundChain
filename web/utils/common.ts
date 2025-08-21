import { assets } from 'chain-registry';
import { Asset, AssetList } from '@chain-registry/types';
import { AminoTypes, GasPrice } from '@cosmjs/stargate';
import { Registry, GeneratedType } from '@cosmjs/proto-signing';
import { SignerOptions, Wallet } from '@cosmos-kit/core';
import { useChain } from '@cosmos-kit/react';

export const getChainAssets = (chainName: string) => {
  return assets.find((chain) => chain.chain_name === chainName) as AssetList;
};

export const getCoin = (chainName: string) => {
  const { assets } = useChain(chainName);
  if (!assets) {
    const chainAssets = getChainAssets(chainName);
    return chainAssets.assets[0] as Asset;
  }

  return assets.assets[0] as Asset;
};

export const getExponent = (chainName: string) => {
  return getCoin(chainName).denom_units.find(
    (unit) => unit.denom === getCoin(chainName).display,
  )?.exponent as number;
};

export const shortenAddress = (address: string, partLength = 6) => {
  return `${address.slice(0, partLength)}...${address.slice(-partLength)}`;
};

export const getWalletLogo = (wallet: Wallet) => {
  if (!wallet?.logo) return '';

  return typeof wallet.logo === 'string'
    ? wallet.logo
    : wallet.logo.major || wallet.logo.minor;
};

// Simple protobuf-like encoding for our custom messages
const encodeString = (str: string): Uint8Array => {
  const bytes = new TextEncoder().encode(str);
  const length = bytes.length;
  const result = new Uint8Array(length + 1);
  result[0] = length; // Length prefix
  result.set(bytes, 1);
  return result;
};

const encodeMessage = (msg: any): Uint8Array => {
  // Simple field encoding: field_number << 3 | wire_type, then data
  const parts: Uint8Array[] = [];
  
  if (msg.creator) {
    parts.push(new Uint8Array([0x0A])); // field 1, wire type 2 (length-delimited)
    parts.push(encodeString(msg.creator));
  }
  if (msg.title) {
    parts.push(new Uint8Array([0x12])); // field 2, wire type 2
    parts.push(encodeString(msg.title));
  }
  if (msg.budget) {
    parts.push(new Uint8Array([0x1A])); // field 3, wire type 2
    parts.push(encodeString(msg.budget));
  }
  if (msg.ipfs_hash) {
    parts.push(new Uint8Array([0x22])); // field 4, wire type 2
    parts.push(encodeString(msg.ipfs_hash));
  }
  if (msg.reviewers && msg.reviewers.length > 0) {
    for (const reviewer of msg.reviewers) {
      parts.push(new Uint8Array([0x2A])); // field 5, wire type 2
      parts.push(encodeString(reviewer));
    }
  }
  if (msg.attest_threshold !== undefined) {
    parts.push(new Uint8Array([0x30])); // field 6, wire type 0 (varint)
    parts.push(new Uint8Array([msg.attest_threshold])); // Simple single byte for small numbers
  }
  if (msg.project_id) {
    parts.push(new Uint8Array([0x12])); // field 2, wire type 2
    parts.push(encodeString(msg.project_id));
  }
  if (msg.milestone_hash) {
    parts.push(new Uint8Array([0x1A])); // field 3, wire type 2
    parts.push(encodeString(msg.milestone_hash));
  }
  
  // Combine all parts
  const totalLength = parts.reduce((sum, part) => sum + part.length, 0);
  const result = new Uint8Array(totalLength);
  let offset = 0;
  for (const part of parts) {
    result.set(part, offset);
    offset += part.length;
  }
  return result;
};

const msgSubmitProjectType: GeneratedType = {
  encode: (message: any) => encodeMessage(message),
  decode: (input: Uint8Array) => ({}), // Not needed for our use case
} as any;

const msgAttestMilestoneType: GeneratedType = {
  encode: (message: any) => encodeMessage(message),
  decode: (input: Uint8Array) => ({}),
} as any;

export const getFundchainRegistry = () => {
  const registry = new Registry();
  registry.register('/fundchain.milestones.v1.MsgSubmitProject', msgSubmitProjectType);
  registry.register('/fundchain.milestones.v1.MsgAttestMilestone', msgAttestMilestoneType);
  // eslint-disable-next-line no-console
  console.log('[fundchain] Proto Registry registered:', ['/fundchain.milestones.v1.MsgSubmitProject', '/fundchain.milestones.v1.MsgAttestMilestone']);
  return registry;
};

export const getFundchainAminoTypes = () => {
  const fundchainAminoConverters: any = {
    '/fundchain.milestones.v1.MsgSubmitProject': {
      aminoType: 'fundchain/MsgSubmitProject',
      toAmino: (msg: any) => ({
        creator: msg.creator,
        title: msg.title,
        budget: msg.budget,
        ipfs_hash: msg.ipfs_hash,
        reviewers: msg.reviewers || [],
        attest_threshold: msg.attest_threshold,
      }),
      fromAmino: (amino: any) => ({
        creator: amino.creator,
        title: amino.title,
        budget: amino.budget,
        ipfs_hash: amino.ipfs_hash,
        reviewers: amino.reviewers || [],
        attest_threshold: Number(amino.attest_threshold || 0),
      }),
    },
    '/fundchain.milestones.v1.MsgAttestMilestone': {
      aminoType: 'fundchain/MsgAttestMilestone',
      toAmino: (msg: any) => ({
        creator: msg.creator,
        project_id: msg.project_id,
        milestone_hash: msg.milestone_hash,
      }),
      fromAmino: (amino: any) => ({
        creator: amino.creator,
        project_id: amino.project_id,
        milestone_hash: amino.milestone_hash,
      }),
    },
  };
  const keys = Object.keys(fundchainAminoConverters || {});
  // eslint-disable-next-line no-console
  console.log('[fundchain] Amino additions registered:', keys);
  return new AminoTypes({ additions: fundchainAminoConverters as any });
};

export const getSignerOptions = (): SignerOptions => {
  const defaultGasPrice = GasPrice.fromString('0.025ufund');

  return {
    // @ts-ignore
    signingStargate: (chain) => {
      if (typeof chain === 'string') {
        return {
          gasPrice: defaultGasPrice,
          aminoTypes: getFundchainAminoTypes(),
          registry: getFundchainRegistry(),
        };
      }
      let gasPrice;
      try {
        const feeToken = chain.fees?.fee_tokens[0];
        const fee = `${feeToken?.average_gas_price || 0.025}${feeToken?.denom || 'ufund'}`;
        gasPrice = GasPrice.fromString(fee);
      } catch (error) {
        gasPrice = defaultGasPrice;
      }
      return {
        gasPrice,
        aminoTypes: getFundchainAminoTypes(),
        registry: getFundchainRegistry(),
      };
    },
    preferredSignType: (chain) => {
      // Force direct signing for fundchain to use our Registry
      if (typeof chain === 'string' && chain === 'fundchain') {
        return 'direct';
      }
      if (typeof chain === 'object' && chain?.chain_name === 'fundchain') {
        return 'direct';
      }
      return 'amino';
    },
  };
};
