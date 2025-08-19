import { AssetList, Chain } from '@chain-registry/types';

export type PageRequest = {
  key?: string;
  offset?: string;
  limit?: string;
  count_total?: boolean;
};

export type ParamsResponse = {
  params: Record<string, unknown>;
};

export type TreasuryBalanceResponse = {
  treasury: string;
  balance: { denom: string; amount: string };
};

export type Project = {
  id: string;
  title: string;
  owner: string;
  total_requested: { denom: string; amount: string };
  status: string;
  reviewers?: string[];
  attest_threshold?: number;
};

export type ProjectsResponse = {
  projects: Project[];
  pagination?: unknown;
};

export type Milestone = {
  id: string;
  project_id: string;
  tranche: number;
  status: string;
  description?: string;
};

export type ProjectMilestonesResponse = {
  milestones: Milestone[];
  pagination?: unknown;
};

export const makeRest = async (
  getRestEndpoint: () => Promise<string | { url?: string; address?: string } | undefined>
) => {
  const baseVal = await getRestEndpoint();
  const base = typeof baseVal === 'string' ? baseVal : baseVal?.url || baseVal?.address;
  if (!base) throw new Error('REST endpoint unavailable');

  const getJson = async <T>(path: string): Promise<T> => {
    const url = `${base}${path}`;
    const res = await fetch(url, { headers: { Accept: 'application/json' } });
    if (!res.ok) {
      let body: string | undefined;
      try {
        body = await res.text();
      } catch {}
      throw new Error(`HTTP ${res.status}: ${path}${body ? `\n${body}` : ''}`);
    }
    return (await res.json()) as T;
  };

  return {
    params: () => getJson<ParamsResponse>('/fundchain/milestones/v1/params'),
    treasuryBalance: () =>
      getJson<TreasuryBalanceResponse>('/fundchain/milestones/v1/treasury/balance'),
    projects: (pagination?: PageRequest) =>
      getJson<ProjectsResponse>(`/fundchain/milestones/v1/projects`),
    project: (id: string) =>
      getJson<{ project: Project }>(`/fundchain/milestones/v1/projects/${id}`),
    projectMilestones: (projectId: string, pagination?: PageRequest) =>
      getJson<ProjectMilestonesResponse>(
        `/fundchain/milestones/v1/projects/${projectId}/milestones`
      ),
  };
};
