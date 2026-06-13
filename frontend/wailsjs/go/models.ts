export namespace main {
	
	export class CalculatedTax {
	    TotalIncomeCents: number;
	    TaxCents: number;
	    ExpensesCents: number;
	    PaidInsuranceCents: number;
	    Quarter: number;
	    Year: number;
	    MonthStart: number;
	    MonthEnd: number;
	
	    static createFrom(source: any = {}) {
	        return new CalculatedTax(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TotalIncomeCents = source["TotalIncomeCents"];
	        this.TaxCents = source["TaxCents"];
	        this.ExpensesCents = source["ExpensesCents"];
	        this.PaidInsuranceCents = source["PaidInsuranceCents"];
	        this.Quarter = source["Quarter"];
	        this.Year = source["Year"];
	        this.MonthStart = source["MonthStart"];
	        this.MonthEnd = source["MonthEnd"];
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
	export class TaxesConfig {
	    MinInsuranceIncomeCents: number;
	    MaxInsuranceIncomeCents: number;
	    ExpensesPercentage: number;
	    TaxPercentage: number;
	    HealthInsurancePercentage: number;
	    PregnancyInsurancePercentage: number;
	    PensionPercentagePartOne: number;
	    PensionPercentagePartTwo: number;
	
	    static createFrom(source: any = {}) {
	        return new TaxesConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MinInsuranceIncomeCents = source["MinInsuranceIncomeCents"];
	        this.MaxInsuranceIncomeCents = source["MaxInsuranceIncomeCents"];
	        this.ExpensesPercentage = source["ExpensesPercentage"];
	        this.TaxPercentage = source["TaxPercentage"];
	        this.HealthInsurancePercentage = source["HealthInsurancePercentage"];
	        this.PregnancyInsurancePercentage = source["PregnancyInsurancePercentage"];
	        this.PensionPercentagePartOne = source["PensionPercentagePartOne"];
	        this.PensionPercentagePartTwo = source["PensionPercentagePartTwo"];
	    }
	}
	export class SocialSecurityParts {
	    PensionPartOneCents: number;
	    PensionPartTwoCents: number;
	    HealthInsuranceCents: number;
	
	    static createFrom(source: any = {}) {
	        return new SocialSecurityParts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PensionPartOneCents = source["PensionPartOneCents"];
	        this.PensionPartTwoCents = source["PensionPartTwoCents"];
	        this.HealthInsuranceCents = source["HealthInsuranceCents"];
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
	    SocialSecurityToPayParts: SocialSecurityParts;
	    TaxesReallyPaidCents: number;
	    SocialSecurityReallyPaidCents: number;
	    SocialSecurityReallyPaidParts: SocialSecurityParts;
	    TaxesConfig: TaxesConfig;
	    Settings: Settings;
	
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
	        this.SocialSecurityToPayParts = this.convertValues(source["SocialSecurityToPayParts"], SocialSecurityParts);
	        this.TaxesReallyPaidCents = source["TaxesReallyPaidCents"];
	        this.SocialSecurityReallyPaidCents = source["SocialSecurityReallyPaidCents"];
	        this.SocialSecurityReallyPaidParts = this.convertValues(source["SocialSecurityReallyPaidParts"], SocialSecurityParts);
	        this.TaxesConfig = this.convertValues(source["TaxesConfig"], TaxesConfig);
	        this.Settings = this.convertValues(source["Settings"], Settings);
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
	
	
	
	export class UserConfig {
	    FirstName: string;
	    MiddleName: string;
	    LastName: string;
	    Egn: string;
	    Bulstat: string;
	    Phone: string;
	    Email: string;
	
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
	        this.Phone = source["Phone"];
	        this.Email = source["Email"];
	    }
	}

}

