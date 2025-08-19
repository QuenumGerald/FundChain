import { assets } from 'chain-registry';
import { Asset, AssetList } from '@chain-registry/types';
import { AminoTypes, GasPrice } from '@cosmjs/stargate';
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

export const getSignerOptions = (): SignerOptions => {
  const defaultGasPrice = GasPrice.fromString('0.025stake');

  return {
    // @ts-ignore
    signingStargate: (chain) => {
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
      if (typeof chain === 'string') {
        return {
          gasPrice: defaultGasPrice,
          aminoTypes: new AminoTypes(fundchainAminoConverters as any),
        };
      }
      let gasPrice;
      try {
        const feeToken = chain.fees?.fee_tokens[0];
        const fee = `${feeToken?.average_gas_price || 0.025}${feeToken?.denom}`;
        gasPrice = GasPrice.fromString(fee);
      } catch (error) {
        gasPrice = defaultGasPrice;
      }
      return {
        gasPrice,
        aminoTypes: new AminoTypes(fundchainAminoConverters as any),
      };
    },
    preferredSignType: () => 'amino',
  };
};
