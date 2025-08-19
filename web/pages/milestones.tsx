import { useEffect, useState } from 'react';
import { Box, Text } from '@interchain-ui/react';
import { Button } from '@/components/common/Button';
import { useQuery } from '@tanstack/react-query';
import { useChain } from '@cosmos-kit/react';

import { useChainStore } from '@/contexts/chain';
import { useConnectChain } from '@/hooks/common/useConnectChain';
import { makeRest, Project } from '@/utils/fundchain';

export default function MilestonesPage() {
  const { selectedChain } = useChainStore();
  const { connect, isWalletConnected, address } = useConnectChain(selectedChain);
  const { getRestEndpoint } = useChain(selectedChain);

  const restFactoryQuery = useQuery({
    queryKey: ['fundchainRest', selectedChain],
    queryFn: () => makeRest(getRestEndpoint),
    enabled: true,
    staleTime: Infinity,
  });

  const rest = restFactoryQuery.data;

  const paramsQuery = useQuery({
    queryKey: ['params', selectedChain],
    queryFn: () => rest!.params(),
    enabled: Boolean(rest),
  });

  const treasuryQuery = useQuery({
    queryKey: ['treasury', selectedChain],
    queryFn: () => rest!.treasuryBalance(),
    enabled: Boolean(rest),
  });

  const projectsQuery = useQuery({
    queryKey: ['projects', selectedChain],
    queryFn: () => rest!.projects(),
    enabled: Boolean(rest),
  });

  const [selectedProjectId, setSelectedProjectId] = useState<string | undefined>(
    undefined
  );

  useEffect(() => {
    const first = projectsQuery.data?.projects?.[0]?.id;
    if (first && !selectedProjectId) setSelectedProjectId(first);
  }, [projectsQuery.data, selectedProjectId]);

  const milestonesQuery = useQuery({
    queryKey: ['milestones', selectedChain, selectedProjectId],
    queryFn: () => rest!.projectMilestones(selectedProjectId!),
    enabled: Boolean(rest) && Boolean(selectedProjectId),
  });

  return (
    <Box display="flex" flexDirection="column" gap="$8" mt="$12">
      <Box as="h2" fontSize="$2xl" fontWeight="$semibold">Milestones Module</Box>

      <Box p="$6" borderRadius="$lg" border="1px solid #e5e7eb" backgroundColor="$cardBg">
        <Box display="flex" alignItems="center" gap="$6">
          <Button onClick={connect} leftIcon="walletFilled" variant="primary">
            {isWalletConnected ? 'Wallet Connected' : 'Connect Wallet'}
          </Button>
          <Text color="$secondaryText">
            {address ? `Address: ${address}` : 'No wallet connected'}
          </Text>
        </Box>
      </Box>

      {/* Submit Project (UI only) */}
      <SubmitProjectForm />

      {/* Attest Milestone (UI only) */}
      <AttestMilestoneForm defaultProjectId={selectedProjectId} />

      <Box p="$6" borderRadius="$lg" border="1px solid #e5e7eb" backgroundColor="$cardBg">
        <Box as="h3" fontSize="$lg" fontWeight="$semibold">Params</Box>
        {paramsQuery.isLoading && <Text>Loading...</Text>}
        {paramsQuery.error && (
          <Text color="$dangerText">{String(paramsQuery.error)}</Text>
        )}
        {paramsQuery.data && (
          <pre style={{ overflowX: 'auto' }}>
            {JSON.stringify(paramsQuery.data.params, null, 2)}
          </pre>
        )}
      </Box>

      <Box p="$6" borderRadius="$lg" border="1px solid #e5e7eb" backgroundColor="$cardBg">
        <Box as="h3" fontSize="$lg" fontWeight="$semibold">Treasury</Box>
        {treasuryQuery.isLoading && <Text>Loading...</Text>}
        {treasuryQuery.error && (
          <Box>
            {String(treasuryQuery.error).includes('treasury or denom parameter not set') ? (
              <Text color="$secondaryText">
                Treasury not configured yet. Set module params (treasury address and denom) to enable this query.
              </Text>
            ) : (
              <Text color="$dangerText">{String(treasuryQuery.error)}</Text>
            )}
          </Box>
        )}
        {treasuryQuery.data && (
          <Box>
            <Text>Address: {treasuryQuery.data.treasury}</Text>
            <Text>
              Balance: {treasuryQuery.data.balance.amount}{' '}
              {treasuryQuery.data.balance.denom}
            </Text>
          </Box>
        )}
      </Box>

      <Box p="$6" borderRadius="$lg" border="1px solid #e5e7eb" backgroundColor="$cardBg">
        <Box as="h3" fontSize="$lg" fontWeight="$semibold">Projects</Box>
        {projectsQuery.isLoading && <Text>Loading...</Text>}
        {projectsQuery.error && (
          <Text color="$dangerText">{String(projectsQuery.error)}</Text>
        )}
        {projectsQuery.data && (
          <Box display="flex" flexDirection="column" gap="$4">
            <select
              value={selectedProjectId}
              onChange={(e) => setSelectedProjectId(e.target.value)}
              style={{ padding: '8px', borderRadius: 8, border: '1px solid #ddd', maxWidth: 360 }}
            >
              {(projectsQuery.data.projects || []).map((p: Project) => (
                <option key={p.id} value={p.id}>{`#${p.id} ${p.title || ''}`}</option>
              ))}
            </select>
            <Box>
              <Box as="h4" fontSize="$md" mb="$2" fontWeight="$semibold">
                Milestones for Project {selectedProjectId}
              </Box>
              {milestonesQuery.isLoading && <Text>Loading...</Text>}
              {milestonesQuery.error && (
                <Text color="$dangerText">{String(milestonesQuery.error)}</Text>
              )}
              {milestonesQuery.data && (
                <pre style={{ overflowX: 'auto' }}>
                  {JSON.stringify(milestonesQuery.data.milestones, null, 2)}
                </pre>
              )}
            </Box>
          </Box>
        )}
      </Box>
    </Box>
  );
}

function SubmitProjectForm() {
  const [title, setTitle] = useState('');
  const [budget, setBudget] = useState('');
  const [ipfs, setIpfs] = useState('');

  const disabled = !title || !budget || !ipfs;

  const onSubmit = () => {
    if (disabled) return;
    // Placeholder: wire to useTx with MsgSubmitProject later
    console.log('SubmitProject payload', { title, budget, ipfs_hash: ipfs });
    alert('SubmitProject UI only. TX wiring pending.');
  };

  return (
    <Box p="$6" borderRadius="$lg" border="1px solid #e5e7eb" backgroundColor="$cardBg">
      <Box as="h3" fontSize="$lg" fontWeight="$semibold">Submit Project</Box>
      <Box display="flex" gap="$4" mt="$3" flexWrap="wrap">
        <input
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          style={{ padding: '8px', borderRadius: 8, border: '1px solid #ddd', minWidth: 220 }}
        />
        <input
          placeholder="Budget (e.g. 1000000ufund)"
          value={budget}
          onChange={(e) => setBudget(e.target.value)}
          style={{ padding: '8px', borderRadius: 8, border: '1px solid #ddd', minWidth: 220 }}
        />
        <input
          placeholder="IPFS Hash"
          value={ipfs}
          onChange={(e) => setIpfs(e.target.value)}
          style={{ padding: '8px', borderRadius: 8, border: '1px solid #ddd', minWidth: 260 }}
        />
        <Button onClick={onSubmit} disabled={disabled} variant="primary">Submit</Button>
      </Box>
    </Box>
  );
}

function AttestMilestoneForm({ defaultProjectId }: { defaultProjectId?: string }) {
  const [projectId, setProjectId] = useState(defaultProjectId || '');
  const [milestoneHash, setMilestoneHash] = useState('');

  // keep in sync if selection above changes
  useEffect(() => {
    if (defaultProjectId && !projectId) setProjectId(defaultProjectId);
  }, [defaultProjectId]);

  const disabled = !projectId || !milestoneHash;

  const onSubmit = () => {
    if (disabled) return;
    // Placeholder: wire to useTx with MsgAttestMilestone later
    console.log('AttestMilestone payload', { project_id: projectId, milestone_hash: milestoneHash });
    alert('AttestMilestone UI only. TX wiring pending.');
  };

  return (
    <Box p="$6" borderRadius="$lg" border="1px solid #e5e7eb" backgroundColor="$cardBg">
      <Box as="h3" fontSize="$lg" fontWeight="$semibold">Attest Milestone</Box>
      <Box display="flex" gap="$4" mt="$3" flexWrap="wrap">
        <input
          placeholder="Project ID"
          value={projectId}
          onChange={(e) => setProjectId(e.target.value)}
          style={{ padding: '8px', borderRadius: 8, border: '1px solid #ddd', minWidth: 160 }}
        />
        <input
          placeholder="Milestone Hash"
          value={milestoneHash}
          onChange={(e) => setMilestoneHash(e.target.value)}
          style={{ padding: '8px', borderRadius: 8, border: '1px solid #ddd', minWidth: 260 }}
        />
        <Button onClick={onSubmit} disabled={disabled} variant="primary">Attest</Button>
      </Box>
    </Box>
  );
}
