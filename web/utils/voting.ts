import dayjs from 'dayjs';
import BigNumber from 'bignumber.js';
import { Chain } from '@chain-registry/types';
import {
  Proposal,
  ProposalStatus,
} from 'interchain-query/cosmos/gov/v1beta1/gov';

export function getChainLogo(chain: Chain) {
  return chain.logo_URIs?.svg || chain.logo_URIs?.png || chain.logo_URIs?.jpeg;
}

export function formatDate(date?: Date) {
  if (!date) return null;
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
}

export function paginate(limit: bigint, reverse: boolean = false) {
  return {
    limit,
    reverse,
    key: new Uint8Array(),
    offset: 0n,
    countTotal: true,
  };
}

export function percent(num: number | string = 0, total: number, decimals = 2) {
  return total
    ? new BigNumber(num)
      .dividedBy(total)
      .multipliedBy(100)
      .decimalPlaces(decimals)
      .toNumber()
    : 0;
}

export const exponentiate = (num: number | string | undefined, exp: number) => {
  if (!num) return 0;
  return new BigNumber(num)
    .multipliedBy(new BigNumber(10).exponentiatedBy(exp))
    .toNumber();
};

export function decodeUint8Array(value?: unknown) {
  if (!value) return '';
  if (typeof value === 'string') return value;
  // Some gql/rest clients may surface bytes as number[]
  if (Array.isArray(value)) {
    try {
      return new TextDecoder('utf-8').decode(Uint8Array.from(value as number[]));
    } catch {
      return '';
    }
  }
  // Normal case: Uint8Array
  if (value instanceof Uint8Array) {
    try {
      return new TextDecoder('utf-8').decode(value);
    } catch {
      return '';
    }
  }
  // Fallback: attempt to coerce objects with a numeric length into Uint8Array
  try {
    const maybe = (value as any) as { length?: number };
    if (typeof maybe?.length === 'number') {
      return new TextDecoder('utf-8').decode(Uint8Array.from(maybe as any));
    }
  } catch {}
  return '';
}

export function getTitle(value?: Uint8Array) {
  return decodeUint8Array(value)
    .slice(0, 250)
    .match(/[A-Z][A-Za-z].*(?=\u0012)/)?.[0];
}

export function parseQuorum(value?: unknown) {
  const quorum = decodeUint8Array(value);
  if (!quorum) return 0;
  // Expect a decimal string like "0.334000000000000000"
  const bn = new BigNumber(quorum);
  if (!bn.isFinite()) return 0;
  return bn.toNumber();
}

export function processProposals(proposals: Proposal[]) {
  const sorted = proposals.sort(
    (a, b) => Number(b.proposalId) - Number(a.proposalId)
  );

  proposals.forEach((proposal) => {
    // @ts-ignore
    if (!proposal.content?.title && proposal.content?.value) {
      // @ts-ignore
      proposal.content.title = getTitle(proposal.content?.value);
    }
  });

  return sorted
    .filter(
      ({ status }) => status === ProposalStatus.PROPOSAL_STATUS_VOTING_PERIOD
    )
    .concat(
      sorted.filter(
        ({ status }) => status !== ProposalStatus.PROPOSAL_STATUS_VOTING_PERIOD
      )
    );
}
