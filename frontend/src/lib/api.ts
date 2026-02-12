// Wails binding call wrappers.
// Calls are routed to Go backend via Wails v2 runtime bindings.

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function getApp(): any {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  return (window as any)["go"]["app"]["App"];
}

export async function bootstrap(): Promise<unknown> {
  return getApp().Bootstrap();
}

export async function saveAppConfig(config: unknown): Promise<void> {
  return getApp().SaveAppConfig(JSON.stringify(config));
}

export async function listProfiles(): Promise<unknown[]> {
  return getApp().ListProfiles();
}

export async function upsertProfile(profile: unknown): Promise<void> {
  return getApp().UpsertProfile(JSON.stringify(profile));
}

export async function deleteProfile(profileId: string): Promise<void> {
  return getApp().DeleteProfile(profileId);
}

export async function duplicateProfile(profileId: string, newName: string): Promise<unknown> {
  return getApp().DuplicateProfile(profileId, newName);
}

export async function setDefaultProfile(profileId: string): Promise<void> {
  return getApp().SetDefaultProfile(profileId);
}

export async function getGPUInfo(): Promise<string> {
  return getApp().GetGPUInfo();
}

export async function detectExternalTools(): Promise<unknown> {
  return getApp().DetectExternalTools();
}

export async function startEncode(request: unknown): Promise<void> {
  return getApp().StartEncode(JSON.stringify(request));
}

export async function requestGracefulStop(sessionId: string): Promise<void> {
  return getApp().RequestGracefulStop(sessionId);
}

export async function requestAbort(sessionId: string): Promise<void> {
  return getApp().RequestAbort(sessionId);
}

export async function cancelJob(sessionId: string, jobId: string): Promise<void> {
  return getApp().CancelJob(sessionId, jobId);
}

export async function resolveOverwrite(sessionId: string, jobId: string, decision: string): Promise<void> {
  return getApp().ResolveOverwrite(sessionId, jobId, decision);
}

export async function listTempArtifacts(): Promise<string[]> {
  return getApp().ListTempArtifacts();
}

export async function cleanupTempArtifacts(paths: string[]): Promise<void> {
  return getApp().CleanupTempArtifacts(paths);
}

export async function getCommandPreview(profile: unknown, inputPath: string, outputPath: string): Promise<string> {
  return getApp().GetCommandPreview(JSON.stringify(profile), inputPath, outputPath);
}
