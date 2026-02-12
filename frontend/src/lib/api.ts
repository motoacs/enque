import type { AppConfig, BootstrapResponse, PreviewCommandResponse, Profile } from './types';

function appApi() {
  const api = window.go?.main?.App;
  if (!api) {
    throw new Error('Wails API is not available');
  }
  return api;
}

export const api = {
  bootstrap: () => appApi().Bootstrap() as Promise<BootstrapResponse>,
  saveAppConfig: (config: AppConfig) => appApi().SaveAppConfig(config),
  listProfiles: () => appApi().ListProfiles() as Promise<Profile[]>,
  upsertProfile: (profile: Profile) => appApi().UpsertProfile(profile) as Promise<Profile>,
  deleteProfile: (profileID: string) => appApi().DeleteProfile(profileID),
  duplicateProfile: (profileID: string, newName: string) => appApi().DuplicateProfile(profileID, newName) as Promise<Profile>,
  setDefaultProfile: (profileID: string) => appApi().SetDefaultProfile(profileID),
  previewCommand: (request: {
    profile: Profile;
    app_config_snapshot: AppConfig;
    input_path: string;
    output_path: string;
  }) => appApi().PreviewCommand(request) as Promise<PreviewCommandResponse>,
  startEncode: (request: unknown) => appApi().StartEncode(request),
  requestGracefulStop: (sessionID: string) => appApi().RequestGracefulStop(sessionID),
  requestAbort: (sessionID: string) => appApi().RequestAbort(sessionID),
  cancelJob: (sessionID: string, jobID: string) => appApi().CancelJob(sessionID, jobID),
  resolveOverwrite: (sessionID: string, jobID: string, decision: 'overwrite' | 'skip' | 'abort') =>
    appApi().ResolveOverwrite(sessionID, jobID, decision),
  detectExternalTools: () => appApi().DetectExternalTools(),
  getGPUInfo: () => appApi().GetGPUInfo(),
  listTempArtifacts: () => appApi().ListTempArtifacts() as Promise<string[]>,
  cleanupTempArtifacts: (paths: string[]) => appApi().CleanupTempArtifacts(paths)
};
