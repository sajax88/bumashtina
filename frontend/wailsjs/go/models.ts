export namespace main {
	
	export class Settings {
	    IsPregnancyInsuranceEnabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IsPregnancyInsuranceEnabled = source["IsPregnancyInsuranceEnabled"];
	    }
	}
	export class TaxesConfig {
	    MinInsuranceIncomeCents: number;
	    MaxInsuranceIncomeCents: number;
	    ExpensesPercentage: number;
	    TaxPercentage: number;
	    PensionPercentage: number;
	    HealthInsurancePercentage: number;
	    PregnancyInsurancePercentage: number;
	
	    static createFrom(source: any = {}) {
	        return new TaxesConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MinInsuranceIncomeCents = source["MinInsuranceIncomeCents"];
	        this.MaxInsuranceIncomeCents = source["MaxInsuranceIncomeCents"];
	        this.ExpensesPercentage = source["ExpensesPercentage"];
	        this.TaxPercentage = source["TaxPercentage"];
	        this.PensionPercentage = source["PensionPercentage"];
	        this.HealthInsurancePercentage = source["HealthInsurancePercentage"];
	        this.PregnancyInsurancePercentage = source["PregnancyInsurancePercentage"];
	    }
	}
	export class UserConfig {
	    FirstName: string;
	    MiddleName: string;
	    LastName: string;
	    Egn: string;
	    Bulstat: string;
	    YearOfBirth: string;
	
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
	        this.YearOfBirth = source["YearOfBirth"];
	    }
	}

}

