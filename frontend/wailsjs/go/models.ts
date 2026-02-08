export namespace main {
	
	export class Config {
	    FirstName: string;
	    MiddleName: string;
	    LastName: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FirstName = source["FirstName"];
	        this.MiddleName = source["MiddleName"];
	        this.LastName = source["LastName"];
	    }
	}

}

