import { Asset, Chain } from '@chain-registry/types';
import { ChainInfo, Currency } from '@keplr-wallet/types';
import { fromBech32 } from '@cosmjs/encoding';

export const makeKeplrChainInfo = (chain: Chain, asset: Asset): ChainInfo => {
  // Resolve a robust decimals value
  const displayUnit = asset?.denom_units?.find(({ denom }) => denom === asset?.display);
  const fallbackUnit = asset?.denom_units?.[asset?.denom_units?.length - 1];
  const decimals = Number.isFinite(displayUnit?.exponent)
    ? (displayUnit!.exponent as number)
    : (Number.isFinite(fallbackUnit?.exponent) ? (fallbackUnit!.exponent as number) : 6);

  const currency: Currency = {
    coinDenom: asset?.symbol || asset?.display || 'FUND',
    coinMinimalDenom: asset?.base || 'ufund',
    coinDecimals: decimals,
    coinGeckoId: asset?.coingecko_id,
    coinImageUrl: asset?.logo_URIs?.svg || asset?.logo_URIs?.png || '',
  };

  return {
    rpc: chain.apis?.rpc?.[0].address ?? 'http://localhost:26657',
    rest: chain.apis?.rest?.[0].address ?? 'http://localhost:1317',
    chainId: chain.chain_id,
    chainName: chain.chain_name,
    bip44: {
      coinType: 118,
    },
    bech32Config: {
      bech32PrefixAccAddr: chain.bech32_prefix ?? '',
      bech32PrefixAccPub: chain.bech32_prefix + 'pub',
      bech32PrefixValAddr: chain.bech32_prefix + 'valoper',
      bech32PrefixValPub: chain.bech32_prefix + 'valoperpub',
      bech32PrefixConsAddr: chain.bech32_prefix + 'valcons',
      bech32PrefixConsPub: chain.bech32_prefix + 'valconspub',
    },
    currencies: [currency],
    feeCurrencies: [
      {
        ...currency,
        gasPriceStep: {
          low: chain.fees?.fee_tokens[0].low_gas_price ?? 0.0025,
          average: chain.fees?.fee_tokens[0].average_gas_price ?? 0.025,
          high: chain.fees?.fee_tokens[0].high_gas_price ?? 0.04,
        },
      },
    ],
    stakeCurrency: currency,
  };
};

export const validateChainAddress = (address: string, bech32Prefix: string) => {
  if (!address.startsWith(bech32Prefix)) {
    return `Invalid prefix (expected "${bech32Prefix}")`;
  }

  try {
    fromBech32(address);
  } catch (e) {
    return 'Invalid address';
  }

  return null;
};
