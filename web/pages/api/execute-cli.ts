import { NextApiRequest, NextApiResponse } from 'next';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { command } = req.body;

  if (!command || typeof command !== 'string') {
    return res.status(400).json({ error: 'Invalid command' });
  }

  // Security: only allow fundchaind commands
  if (!command.startsWith('fundchaind ')) {
    return res.status(400).json({ error: 'Only fundchaind commands allowed' });
  }

  try {
    const { stdout, stderr } = await execAsync(command, {
      cwd: '/home/nova/Documents/projects/Blockchain/FundChain',
      timeout: 30000, // 30 second timeout
    });

    res.status(200).json({
      success: true,
      stdout,
      stderr,
    });
  } catch (error: any) {
    console.error('CLI execution error:', error);
    res.status(500).json({
      success: false,
      error: error.message,
      stdout: error.stdout,
      stderr: error.stderr,
    });
  }
}
