/// <reference types="vite/client" />

import type { AppConfig, BootstrapResponse, PreviewCommandResponse, Profile } from './lib/types';

declare global {
  interface Window {
    go?: {
      main?: {
        App?: {
          Bootstrap: () => Promise<BootstrapResponse>;
          SaveAppConfig: (config: AppConfig) => Promise<void>;
          ListProfiles: () => Promise<Profile[]>;
          UpsertProfile: (profile: Profile) => Promise<Profile>;
          DeleteProfile: (profileID: string) => Promise<void>;
          DuplicateProfile: (profileID: string, newName: string) => Promise<Profile>;
          SetDefaultProfile: (profileID: string) => Promise<void>;
          GetGPUInfo: () => Promise<unknown>;
          DetectExternalTools: () => Promise<unknown>;
          StartEncode: (request: unknown) => Promise<unknown>;
          RequestGracefulStop: (sessionID: string) => Promise<void>;
          RequestAbort: (sessionID: string) => Promise<void>;
          CancelJob: (sessionID: string, jobID: string) => Promise<void>;
          ResolveOverwrite: (sessionID: string, jobID: string, decision: string) => Promise<void>;
          ListTempArtifacts: () => Promise<string[]>;
          CleanupTempArtifacts: (paths: string[]) => Promise<void>;
          PreviewCommand: (request: unknown) => Promise<PreviewCommandResponse>;
        };
      };
    };
    runtime?: {
      EventsOn: (name: string, cb: (payload: unknown) => void) => () => void;
    };
  }
}

export {};
