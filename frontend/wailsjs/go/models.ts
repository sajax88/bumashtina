export namespace main {
	
	export class TaxesConfig {
	    MinInsuranceIncomeCents: number;
	    MaxInsuranceIncomeCents: number;
	    ExpensesPercentage: number;
	    TaxPercentage: number;
	    PensionPercentage: number;
	    HealthInsurancePercentage: number;
	    PregnancyInsurancePercentage: number;
	    PensionPercentagePartOne: number;
	    PensionPercentagePartTwo: number;
	    Divider: number;
	
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
	        this.PensionPercentagePartOne = source["PensionPercentagePartOne"];
	        this.PensionPercentagePartTwo = source["PensionPercentagePartTwo"];
	        this.Divider = source["Divider"];
	    }
	}
	export class IncomeForm {
	    Month: number;
	    Year: number;
	    MonthIncomeCents: number;
	    TaxedIncomeCents: number;
	    DayEnd: number;
	    DayStart: number;
	    WorkDaysTotal: number;
	    WorkDaysReal: number;
	    WorkDaysSickLeave: number;
	    TaxesToPayCents: number;
	    SocialSecurityToPayCents: number;
	    TaxesConfig: TaxesConfig;
	
	    static createFrom(source: any = {}) {
	        return new IncomeForm(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Month = source["Month"];
	        this.Year = source["Year"];
	        this.MonthIncomeCents = source["MonthIncomeCents"];
	        this.TaxedIncomeCents = source["TaxedIncomeCents"];
	        this.DayEnd = source["DayEnd"];
	        this.DayStart = source["DayStart"];
	        this.WorkDaysTotal = source["WorkDaysTotal"];
	        this.WorkDaysReal = source["WorkDaysReal"];
	        this.WorkDaysSickLeave = source["WorkDaysSickLeave"];
	        this.TaxesToPayCents = source["TaxesToPayCents"];
	        this.SocialSecurityToPayCents = source["SocialSecurityToPayCents"];
	        this.TaxesConfig = this.convertValues(source["TaxesConfig"], TaxesConfig);
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

