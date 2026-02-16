export namespace main {
	
	export class UserConfig {
	    FirstName: string;
	    MiddleName: string;
	    LastName: string;
	    Egn: string;
	    Bulstat: string;
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FirstName = source["FirstName"];
	        this.MiddleName = source["MiddleName"];
	        this.LastName = source["LastName"];
	        this.Egn = source["Egn"];
	        this.Bulstat = source["Bulstat"];
	    }
	}

}

