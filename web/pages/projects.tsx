import { useEffect, useMemo, useState } from 'react';
import { Box, Text } from '@interchain-ui/react';
import { useQuery } from '@tanstack/react-query';
import { useChain } from '@cosmos-kit/react';

import { useChainStore } from '@/contexts/chain';
import { makeRest, Project, Milestone } from '@/utils/fundchain';

export default function ProjectsPage() {
  const { selectedChain } = useChainStore();
  const { getRestEndpoint } = useChain(selectedChain);

  const restFactoryQuery = useQuery({
    queryKey: ['fundchainRest', selectedChain],
    queryFn: () => makeRest(getRestEndpoint),
    enabled: true,
    staleTime: Infinity,
  });
  const rest = restFactoryQuery.data;

  const projectsQuery = useQuery({
    queryKey: ['projects', selectedChain],
    queryFn: () => rest!.projects(),
    enabled: Boolean(rest),
  });

  const [expandedProjectId, setExpandedProjectId] = useState<string | null>(null);

  const milestonesQueries = useMemo(() => {
    if (!rest || !projectsQuery.data) return {} as Record<string, ReturnType<typeof useQuery<unknown>>>;
    // We'll fetch milestones when a project card is expanded
    return {} as Record<string, ReturnType<typeof useQuery<unknown>>>;
  }, [rest, projectsQuery.data]);

  return (
    <Box display="flex" flexDirection="column" gap="$8" mt="$12">
      <Box as="h2" fontSize="$2xl" fontWeight="$semibold">Projects Progress</Box>

      {projectsQuery.isLoading && <Text>Loading projects...</Text>}
      {projectsQuery.error && (
        <Text color="$dangerText">{String(projectsQuery.error)}</Text>
      )}

      {projectsQuery.data && (
        <Box display="flex" flexDirection="column" gap="$4">
          {(projectsQuery.data.projects || []).map((p: Project) => (
            <ProjectCard
              key={p.id}
              project={p}
              restReady={Boolean(rest)}
              selectedChain={selectedChain}
              expanded={expandedProjectId === p.id}
              onToggle={() => setExpandedProjectId(expandedProjectId === p.id ? null : p.id)}
            />
          ))}
        </Box>
      )}
    </Box>
  );
}

function ProjectCard({
  project,
  restReady,
  selectedChain,
  expanded,
  onToggle,
}: {
  project: Project;
  restReady: boolean;
  selectedChain: string;
  expanded: boolean;
  onToggle: () => void;
}) {
  const { getRestEndpoint } = useChain(selectedChain);
  const restFactoryQuery = useQuery({
    queryKey: ['fundchainRest', selectedChain],
    queryFn: () => makeRest(getRestEndpoint),
    enabled: true,
    staleTime: Infinity,
  });
  const rest = restFactoryQuery.data;

  const milestonesQuery = useQuery({
    queryKey: ['milestones', selectedChain, project.id],
    queryFn: () => rest!.projectMilestones(project.id),
    enabled: Boolean(rest) && expanded,
  });

  const progress = useMemo(() => {
    const ms: Milestone[] = (milestonesQuery.data?.milestones as Milestone[]) || [];
    const total = ms.length;
    if (total === 0) return 0;
    const completed = ms.filter((m) =>
      typeof m.status === 'string' && m.status.toLowerCase().includes('attest') || m.status.toLowerCase().includes('complete')
    ).length;
    return Math.round((completed / total) * 100);
  }, [milestonesQuery.data]);

  return (
    <Box p="$6" borderRadius="$lg" border="1px solid #e5e7eb" backgroundColor="$cardBg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb="$3">
        <Box>
          <Box as="h3" fontSize="$lg" fontWeight="$semibold">{`#${project.id} ${project.title}`}</Box>
          <Text color="$secondaryText">Owner: {project.owner}</Text>
          <Text color="$secondaryText">
            Budget: {project.total_requested?.amount} {project.total_requested?.denom}
          </Text>
        </Box>
        <button
          onClick={onToggle}
          style={{ padding: '8px 12px', borderRadius: 8, border: '1px solid #ddd', cursor: 'pointer' }}
        >
          {expanded ? 'Hide' : 'View'}
        </button>
      </Box>

      {expanded && (
        <Box>
          {milestonesQuery.isLoading && <Text>Loading milestones...</Text>}
          {milestonesQuery.error && (
            <Text color="$dangerText">{String(milestonesQuery.error)}</Text>
          )}
          {milestonesQuery.data && (
            <Box>
              <Box mb="$2">Progress</Box>
              <ProgressBar percent={progress} />
              <Box mt="$3">
                <pre style={{ overflowX: 'auto' }}>
                  {JSON.stringify(milestonesQuery.data.milestones, null, 2)}
                </pre>
              </Box>
            </Box>
          )}
        </Box>
      )}
    </Box>
  );
}

function ProgressBar({ percent }: { percent: number }) {
  return (
    <Box width="100%" height="12px" borderRadius="8px" backgroundColor="#eee">
      <Box
        height="12px"
        width={`${Math.min(Math.max(percent, 0), 100)}%`}
        backgroundColor="#7c3aed"
        borderRadius="8px"
      />
    </Box>
  );
}
