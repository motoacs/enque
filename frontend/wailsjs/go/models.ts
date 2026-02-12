export namespace app {
	
	export class BootstrapResult {
	    config: config.AppConfig;
	    profiles: profile.Profile[];
	    tools: detector.DetectionResult;
	    temp_artifacts: string[];
	
	    static createFrom(source: any = {}) {
	        return new BootstrapResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config = this.convertValues(source["config"], config.AppConfig);
	        this.profiles = this.convertValues(source["profiles"], profile.Profile);
	        this.tools = this.convertValues(source["tools"], detector.DetectionResult);
	        this.temp_artifacts = source["temp_artifacts"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace config {
	
	export class AppConfig {
	    version: number;
	    nvencc_path: string;
	    qsvenc_path: string;
	    ffmpeg_path: string;
	    ffprobe_path: string;
	    max_concurrent_jobs: number;
	    on_error: string;
	    decoder_fallback: boolean;
	    keep_failed_temp: boolean;
	    no_output_timeout_sec: number;
	    no_progress_timeout_sec: number;
	    post_complete_action: string;
	    post_complete_command: string;
	    output_folder_mode: string;
	    output_folder_path: string;
	    output_name_template: string;
	    output_container: string;
	    overwrite_mode: string;
	    language: string;
	    default_profile_id: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.nvencc_path = source["nvencc_path"];
	        this.qsvenc_path = source["qsvenc_path"];
	        this.ffmpeg_path = source["ffmpeg_path"];
	        this.ffprobe_path = source["ffprobe_path"];
	        this.max_concurrent_jobs = source["max_concurrent_jobs"];
	        this.on_error = source["on_error"];
	        this.decoder_fallback = source["decoder_fallback"];
	        this.keep_failed_temp = source["keep_failed_temp"];
	        this.no_output_timeout_sec = source["no_output_timeout_sec"];
	        this.no_progress_timeout_sec = source["no_progress_timeout_sec"];
	        this.post_complete_action = source["post_complete_action"];
	        this.post_complete_command = source["post_complete_command"];
	        this.output_folder_mode = source["output_folder_mode"];
	        this.output_folder_path = source["output_folder_path"];
	        this.output_name_template = source["output_name_template"];
	        this.output_container = source["output_container"];
	        this.overwrite_mode = source["overwrite_mode"];
	        this.language = source["language"];
	        this.default_profile_id = source["default_profile_id"];
	    }
	}

}

export namespace detector {
	
	export class ToolInfo {
	    name: string;
	    path: string;
	    version: string;
	    found: boolean;
	    error?: string;
	    supported: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ToolInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.version = source["version"];
	        this.found = source["found"];
	        this.error = source["error"];
	        this.supported = source["supported"];
	    }
	}
	export class DetectionResult {
	    nvencc: ToolInfo;
	    qsvenc: ToolInfo;
	    ffmpeg: ToolInfo;
	    ffprobe: ToolInfo;
	
	    static createFrom(source: any = {}) {
	        return new DetectionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nvencc = this.convertValues(source["nvencc"], ToolInfo);
	        this.qsvenc = this.convertValues(source["qsvenc"], ToolInfo);
	        this.ffmpeg = this.convertValues(source["ffmpeg"], ToolInfo);
	        this.ffprobe = this.convertValues(source["ffprobe"], ToolInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace profile {
	
	export class NVEncCAdvanced {
	    interlace: string;
	    avsw_decoder: string;
	    input_csp: string;
	    output_csp: string;
	    tune: string;
	    max_bitrate?: number;
	    vbr_quality?: number;
	    lookahead_level?: number;
	    weightp: boolean;
	    mv_precision: string;
	    refs_forward?: number;
	    refs_backward?: number;
	    level: string;
	    profile: string;
	    tier: string;
	    ssim: boolean;
	    psnr: boolean;
	    trim: string;
	    seek: string;
	    seekto: string;
	    video_metadata: string;
	    audio_copy: string;
	    audio_codec: string;
	    audio_bitrate: string;
	    audio_quality: string;
	    audio_samplerate: string;
	    audio_metadata: string;
	    sub_copy: string;
	    sub_metadata: string;
	    data_copy: string;
	    attachment_copy: string;
	    metadata: string;
	    output_thread?: number;
	
	    static createFrom(source: any = {}) {
	        return new NVEncCAdvanced(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.interlace = source["interlace"];
	        this.avsw_decoder = source["avsw_decoder"];
	        this.input_csp = source["input_csp"];
	        this.output_csp = source["output_csp"];
	        this.tune = source["tune"];
	        this.max_bitrate = source["max_bitrate"];
	        this.vbr_quality = source["vbr_quality"];
	        this.lookahead_level = source["lookahead_level"];
	        this.weightp = source["weightp"];
	        this.mv_precision = source["mv_precision"];
	        this.refs_forward = source["refs_forward"];
	        this.refs_backward = source["refs_backward"];
	        this.level = source["level"];
	        this.profile = source["profile"];
	        this.tier = source["tier"];
	        this.ssim = source["ssim"];
	        this.psnr = source["psnr"];
	        this.trim = source["trim"];
	        this.seek = source["seek"];
	        this.seekto = source["seekto"];
	        this.video_metadata = source["video_metadata"];
	        this.audio_copy = source["audio_copy"];
	        this.audio_codec = source["audio_codec"];
	        this.audio_bitrate = source["audio_bitrate"];
	        this.audio_quality = source["audio_quality"];
	        this.audio_samplerate = source["audio_samplerate"];
	        this.audio_metadata = source["audio_metadata"];
	        this.sub_copy = source["sub_copy"];
	        this.sub_metadata = source["sub_metadata"];
	        this.data_copy = source["data_copy"];
	        this.attachment_copy = source["attachment_copy"];
	        this.metadata = source["metadata"];
	        this.output_thread = source["output_thread"];
	    }
	}
	export class Profile {
	    id: string;
	    version: number;
	    name: string;
	    is_preset: boolean;
	    encoder_type: string;
	    encoder_options: Record<string, any>;
	    codec: string;
	    rate_control: string;
	    rate_value: number;
	    preset: string;
	    output_depth: number;
	    multipass: string;
	    output_res: string;
	    bframes?: number;
	    ref?: number;
	    lookahead?: number;
	    gop_len?: number;
	    aq: boolean;
	    aq_temporal: boolean;
	    split_enc: string;
	    parallel: string;
	    decoder: string;
	    device: string;
	    audio_mode: string;
	    audio_bitrate: number;
	    colormatrix: string;
	    transfer: string;
	    colorprim: string;
	    colorrange: string;
	    dhdr10_info: string;
	    metadata_copy: boolean;
	    video_metadata_copy: boolean;
	    audio_metadata_copy: boolean;
	    chapter_copy: boolean;
	    sub_copy: boolean;
	    data_copy: boolean;
	    attachment_copy: boolean;
	    restore_file_time: boolean;
	    nvencc_advanced: NVEncCAdvanced;
	    output_container: string;
	    custom_options: string;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.version = source["version"];
	        this.name = source["name"];
	        this.is_preset = source["is_preset"];
	        this.encoder_type = source["encoder_type"];
	        this.encoder_options = source["encoder_options"];
	        this.codec = source["codec"];
	        this.rate_control = source["rate_control"];
	        this.rate_value = source["rate_value"];
	        this.preset = source["preset"];
	        this.output_depth = source["output_depth"];
	        this.multipass = source["multipass"];
	        this.output_res = source["output_res"];
	        this.bframes = source["bframes"];
	        this.ref = source["ref"];
	        this.lookahead = source["lookahead"];
	        this.gop_len = source["gop_len"];
	        this.aq = source["aq"];
	        this.aq_temporal = source["aq_temporal"];
	        this.split_enc = source["split_enc"];
	        this.parallel = source["parallel"];
	        this.decoder = source["decoder"];
	        this.device = source["device"];
	        this.audio_mode = source["audio_mode"];
	        this.audio_bitrate = source["audio_bitrate"];
	        this.colormatrix = source["colormatrix"];
	        this.transfer = source["transfer"];
	        this.colorprim = source["colorprim"];
	        this.colorrange = source["colorrange"];
	        this.dhdr10_info = source["dhdr10_info"];
	        this.metadata_copy = source["metadata_copy"];
	        this.video_metadata_copy = source["video_metadata_copy"];
	        this.audio_metadata_copy = source["audio_metadata_copy"];
	        this.chapter_copy = source["chapter_copy"];
	        this.sub_copy = source["sub_copy"];
	        this.data_copy = source["data_copy"];
	        this.attachment_copy = source["attachment_copy"];
	        this.restore_file_time = source["restore_file_time"];
	        this.nvencc_advanced = this.convertValues(source["nvencc_advanced"], NVEncCAdvanced);
	        this.output_container = source["output_container"];
	        this.custom_options = source["custom_options"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

