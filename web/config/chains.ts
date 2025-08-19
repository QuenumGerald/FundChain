import { AssetList, Chain } from '@chain-registry/types';
import { chains } from 'chain-registry';

// Local FundChain definition (from project's chain.json)
const fundchain: Chain = {
  chain_name: 'fundchain',
  status: 'live',
  network_type: 'devnet',
  pretty_name: 'FundChain',
  chain_type: 'cosmos',
  chain_id: 'fundchain',
  bech32_prefix: 'fund',
  daemon_name: 'fundchaind',
  node_home: '~/.fundchain',
  slip44: 118,
  apis: {
    rpc: [{ address: 'http://localhost:26657', provider: 'localhost' }],
    rest: [{ address: 'http://localhost:1317', provider: 'localhost' }],
    grpc: [{ address: 'localhost:9090', provider: 'localhost' }],
  },
};

export const fundchainAssets: AssetList = {
  chain_name: 'fundchain',
  assets: [
    {
      description: 'Staking token for FundChain devnet',
      denom_units: [
        { denom: 'stake', exponent: 0 },
        { denom: 'STAKE', exponent: 6 },
      ],
      type_asset: 'sdk.coin',
      base: 'stake',
      name: 'Stake',
      display: 'STAKE',
      symbol: 'STAKE',
    },
  ],
};

const chainNames = ['fundchain', 'cosmoshub'];

export const chainOptions = chainNames
  .map((chainName) =>
    chainName === 'fundchain'
      ? fundchain
      : chains.find((chain) => chain.chain_name === chainName)!
  )
  .filter(Boolean);

export { fundchain };