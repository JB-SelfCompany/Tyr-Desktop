export namespace main {
	
	export class BackupOptionsDTO {
	    includeDatabase: boolean;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new BackupOptionsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.includeDatabase = source["includeDatabase"];
	        this.password = source["password"];
	    }
	}
	export class PeerConfigDTO {
	    address: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PeerConfigDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.enabled = source["enabled"];
	    }
	}
	export class ConfigDTO {
	    onboardingComplete: boolean;
	    peers: PeerConfigDTO[];
	    language: string;
	    theme: string;
	    autoStart: boolean;
	    smtpAddress: string;
	    imapAddress: string;
	
	    static createFrom(source: any = {}) {
	        return new ConfigDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.onboardingComplete = source["onboardingComplete"];
	        this.peers = this.convertValues(source["peers"], PeerConfigDTO);
	        this.language = source["language"];
	        this.theme = source["theme"];
	        this.autoStart = source["autoStart"];
	        this.smtpAddress = source["smtpAddress"];
	        this.imapAddress = source["imapAddress"];
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
	
	export class PeerInfoDTO {
	    address: string;
	    enabled: boolean;
	    connected: boolean;
	    latency: number;
	    uptime: number;
	    rxBytes: number;
	    txBytes: number;
	    rxRate: number;
	    txRate: number;
	    lastError?: string;
	
	    static createFrom(source: any = {}) {
	        return new PeerInfoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.enabled = source["enabled"];
	        this.connected = source["connected"];
	        this.latency = source["latency"];
	        this.uptime = source["uptime"];
	        this.rxBytes = source["rxBytes"];
	        this.txBytes = source["txBytes"];
	        this.rxRate = source["rxRate"];
	        this.txRate = source["txRate"];
	        this.lastError = source["lastError"];
	    }
	}
	export class RestoreOptionsDTO {
	    backupPath: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new RestoreOptionsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.backupPath = source["backupPath"];
	        this.password = source["password"];
	    }
	}
	export class ResultDTO {
	    success: boolean;
	    message?: string;
	    data?: string;
	
	    static createFrom(source: any = {}) {
	        return new ResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = source["data"];
	    }
	}
	export class ServiceStatusDTO {
	    status: string;
	    running: boolean;
	    mailAddress: string;
	    smtpAddress: string;
	    imapAddress: string;
	    databasePath: string;
	    errorMessage?: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceStatusDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.running = source["running"];
	        this.mailAddress = source["mailAddress"];
	        this.smtpAddress = source["smtpAddress"];
	        this.imapAddress = source["imapAddress"];
	        this.databasePath = source["databasePath"];
	        this.errorMessage = source["errorMessage"];
	    }
	}

}

