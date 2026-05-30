<script lang="ts">
    import {
        ArrowBigDown,
        ArrowBigUp,
        BookText,
        Calculator,
        Check,
        CircleCheck,
        Plus,
        Save
    } from 'lucide-svelte';
    import {
        CalculateTaxForQuarter,
        GenerateDeclarationOne,
        GenerateDeclarationSix,
        LoadTaxesConfig,
        LoadThisMonthActions,
        SaveIncomeForm,
        SavePaidTaxForQuarter
    } from "../../wailsjs/go/main/App.js";
    import {BrowserOpenURL} from "../../wailsjs/runtime";
    import {onMount} from 'svelte';

    import Info from "../components/Info.svelte"
    import ActionsBanner from "../components/ActionsBanner.svelte";

    let form = {
        Month: String(new Date().getMonth()), // We want a previous month
        Year: new Date().getFullYear(),
        MonthIncome: "",
        TaxedIncome: "",
        DayStart: 0,
        DayEnd: 0,
        WorkDaysTotal: 0,
        WorkDaysSickLeave: 0,
    }


    // Todo: move all this into a separate component
    let previousQuarter = Math.floor((new Date().getMonth() + 3) / 3) - 1;
    let taxFormYear = new Date().getFullYear();
    if (previousQuarter < 1) {
        previousQuarter = 4;
        taxFormYear -= 1
    }

    let taxCalculatorForm = {
        Quarter: previousQuarter.toString(),
        Year: taxFormYear,
    }

    let taxEnterForm = {
        Quarter: previousQuarter.toString(),
        Year: taxFormYear,
        AmountPaid: 0.0,
    }

    let declarationSixForm = {
        Year: new Date().getFullYear(),
    }
    let configTaxes

    let res = "" // TODO
    let data // TODO
    let taxCalculationResult;

    const MONEY_DIVIDER = 100;

    onMount(async () => {
        await load_configs()
    });

    async function load_configs(): Promise<void> {
        const [taxes] = await Promise.all([
            LoadTaxesConfig(),
        ]);
        configTaxes = taxes;
    }

    function setMinIncome(): void {
        form.TaxedIncome = String(configTaxes.MinInsuranceIncomeCents / MONEY_DIVIDER);
    }

    function setMaxIncome(): void {
        form.TaxedIncome = String(configTaxes.MaxInsuranceIncomeCents / MONEY_DIVIDER);
    }

    function saveIncome(): void {
        if (form.WorkDaysSickLeave > form.WorkDaysTotal) {
            alert("Дните в болничен не могат да надвишават общите работни дни.");
            return;
        }

        if (form.DayStart > 0 && form.DayEnd > 0 && form.DayStart >= form.DayEnd) {
            alert("Началният ден трябва да бъде преди крайния ден.");
            return;
        }

        if (form.MonthIncome === "" || form.TaxedIncome === "" || form.Year === 0) {
            alert("Моля, попълнете всички задължителни полета.");
            return;
        }

        if (
            parseFloat(form.TaxedIncome) * MONEY_DIVIDER < configTaxes.MinInsuranceIncomeCents
            || parseFloat(form.TaxedIncome) * MONEY_DIVIDER > configTaxes.MaxInsuranceIncomeCents
        ) {
            alert(`Осигурителният доход трябва да бъде между ${configTaxes.MinInsuranceIncomeCents / MONEY_DIVIDER} и ${configTaxes.MaxInsuranceIncomeCents / MONEY_DIVIDER}.`);
            return;
        }

        // TODO: validation

        let formToSave = {
            Month: parseInt(form.Month),
            Year: form.Year,
            MonthIncomeCents: Math.ceil(parseFloat(form.MonthIncome) * MONEY_DIVIDER),
            TaxedIncomeCents: Math.ceil(parseFloat(form.TaxedIncome) * MONEY_DIVIDER),
            DayStart: form.DayStart,
            DayEnd: form.DayEnd,
            WorkDaysTotal: form.WorkDaysTotal,
            WorkDaysSickLeave: form.WorkDaysSickLeave,
        }

        SaveIncomeForm(formToSave).then((result) => {
            data = result; // TODO: show success message, or error.
            // TODO Hide after a few seconds
        });
    }

    function generateDeclarationOne(): void {
        GenerateDeclarationOne(parseInt(form.Month), form.Year).then((result) => (res = result));
    }

    function generateDeclarationSix(): void {
        GenerateDeclarationSix(declarationSixForm.Year).then((result) => (res = result));
    }

    function calculateTaxForQuarter() {
        CalculateTaxForQuarter(parseInt(taxCalculatorForm.Quarter), taxCalculatorForm.Year).then(
            function (result) {
                taxCalculationResult = result
                taxEnterForm.AmountPaid = result.TaxCents / MONEY_DIVIDER;
                taxEnterForm.Quarter = taxCalculatorForm.Quarter;
                taxEnterForm.Year = taxCalculatorForm.Year;
            }
        );
    }

    let taxEnterResult = "";

    function savePaidTaxForQuarter(): void {
        SavePaidTaxForQuarter(
            parseInt(taxEnterForm.Quarter),
            taxEnterForm.Year,
            parseFloat(taxEnterForm.AmountPaid)
        ).then((result) => (taxEnterResult = result)); // TODO: result
    }
</script>

<main>

    <ActionsBanner />

    <div id="home-page-block-right" class="input-box">

        <!-- TODO: this and tax - to Data page? -->
        <div id="declaration-six-box">
            <div class="form-row">
                <div class="form-group">
                    <button class="btn"
                            on:click={() => {document.getElementById('declaration-six-block').style.display = 'block';}}>
                        <span><BookText color="#444" size="20"/> Генерирай Декларация 6</span>
                    </button>
                    <div id="declaration-six-block" class="hidden-form-block" style="display: none;">
                        <input type="number" id="tax-calculator-year" bind:value={declarationSixForm.Year}/>

                        <button class="btn btn-small" on:click={generateDeclarationSix}>
                            <span><Check color="#444" size="20"/></span>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div id="tax-calculator-box">
            <div class="form-row">
                <div class="form-group">
                    <button class="btn"
                            on:click={() => {document.getElementById('tax-calculator-block').style.display = 'block';}}>
                        <span><Calculator color="#444" size="20"/> Изчисли данък за тримесечие</span>
                    </button>
                    <div id="tax-calculator-block" class="hidden-form-block" style="display: none;">
                        <select id="tax-calculator-quarter" bind:value={taxCalculatorForm.Quarter}>
                            <option value="1">I</option>
                            <option value="2">II</option>
                            <option value="3">III</option>
                            <option value="4">IV</option>
                        </select>

                        <input type="number" id="tax-calculator-year" bind:value={taxCalculatorForm.Year}/>

                        <button class="btn btn-small" on:click={() => {calculateTaxForQuarter()}}>
                            <span><Check color="#444" size="20"/></span>
                        </button>
                        {#if taxCalculationResult}
                            <div class="alert alert-info" id="tax-calculator-result">
                                Месеци: {taxCalculationResult.MonthStart}-{taxCalculationResult.MonthEnd}<br>
                                Доход: {taxCalculationResult.TotalIncomeCents / MONEY_DIVIDER} EUR<br>
                                Приспадащи се разходи: {taxCalculationResult.ExpensesCents / MONEY_DIVIDER} EUR<br>
                                Осигуровки: {taxCalculationResult.PaidInsuranceCents / MONEY_DIVIDER} EUR<br>
                                Данък: {taxCalculationResult.TaxCents / MONEY_DIVIDER} EUR
                                <!-- TODO: full formula -->
                            </div>
                            <!-- TODO: warn if not entered
                            <small><i>Точният данък зависи от фактически платени осигуровки</i></small>
                            -->
                        {/if}
                    </div>
                </div>
            </div>
        </div>

        <div id="tax-calculator-box">
            <div class="form-row">
                <div class="form-group">
                    <button class="btn"
                            on:click={() => {document.getElementById('tax-enter-form').style.display = 'block';}}>
                        <span><Plus color="#444" size="20"/> Въведи платен данък</span>
                    </button>
                    <div id="tax-enter-form" class="hidden-form-block" style="display: none;">
                        <select id="tax-enter-form-quarter" bind:value={taxEnterForm.Quarter}>
                            <option value="1">I</option>
                            <option value="2">II</option>
                            <option value="3">III</option>
                            <option value="4">IV</option>
                        </select>

                        <input type="number" id="tax-enter-form-year" bind:value={taxEnterForm.Year}/>

                        <input type="text" id="tax-enter-form-amount" bind:value={taxEnterForm.AmountPaid}/> EUR

                        <button class="btn btn-small" on:click={() => {savePaidTaxForQuarter()}}>
                            <span><Check color="#444" size="20"/></span>
                        </button>

                        {#if taxEnterResult}
                            <div class="alert success" id="tax-enter-form-result">
                                {taxEnterResult}
                            </div>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div class="input-box" id="home-page-input-box">

        <h2>Въведи данни за доходи</h2>

        <div class="form-row">
            <div class="form-group">
                <label for="Month">Месец</label>
                <select class="input" required id="Month" bind:value={form.Month}>
                    <option value="1">Януари</option>
                    <option value="2">Февруари</option>
                    <option value="3">Март</option>
                    <option value="4">Април</option>
                    <option value="5">Май</option>
                    <option value="6">Юни</option>
                    <option value="7">Юли</option>
                    <option value="8">Август</option>
                    <option value="9">Септември</option>
                    <option value="10">Октомври</option>
                    <option value="11">Ноември</option>
                    <option value="12">Декември</option>
                </select>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="Year">Година</label>
                <input class="input" min="2026" required id="Year" type="number" bind:value={form.Year}/>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="MonthIncome">Доход за месец</label>
                <input class="input" min="0" required id="MonthIncome" type="text" bind:value={form.MonthIncome}/>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="TaxedIncome">Осигурителен доход</label>
                <input class="input" min="0" required id="TaxedIncome" type="text" bind:value={form.TaxedIncome}/>

                <button class="btn btn-small" on:click={setMinIncome}>
                    <span><ArrowBigDown color="#444" size="16"/> Мин</span>
                </button>

                <button class="btn btn-small" on:click={setMaxIncome}>
                    <span><ArrowBigUp color="#444" size="16"/> Макс</span>
                </button>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="WorkDaysTotal">Общо работни дни</label>
                <input class="input" min="0" max="31" id="WorkDaysTotal" type="number" bind:value={form.WorkDaysTotal}/>

                <div class="info">Провери работни дни: <a href=""
                                                          on:click={() => BrowserOpenURL('https://kik-info.com/spravochnik/calendar/' + form.Year)}
                >
                    https://kik-info.com/spravochnik/calendar/{form.Year}
                </a></div>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group checkbox-group">
                <input id="ShowSickLeave" type="checkbox" on:change={(e) => {
                    document.getElementById('SickLeaveFieldset').style.display = e.target.checked ? 'block' : 'none';
                    if (!e.target.checked) {
                      form.WorkDaysSickLeave = 0;
                    }
                }}/>

                <label class='checkbox-label' for="ShowSickLeave">Бях в болничен</label>
            </div>
        </div>

        <fieldset id="SickLeaveFieldset" style="display: none;">
            <div class="form-row">
                <div class="form-group">
                    <label for="WorkDaysSickLeave">Дни в болничен</label>
                    <input class="input" min="0" max="31" id="WorkDaysSickLeave" type="number"
                           bind:value={form.WorkDaysSickLeave}/>
                </div>
            </div>
        </fieldset>

        <div class="form-row">
            <div class="form-group checkbox-group">
                <input id="AddStartOrEnd" type="checkbox" on:change={(e) => {
        document.getElementById('AddStartOrEndFieldset').style.display = e.target.checked ? 'block' : 'none';
        if (!e.target.checked) {
          form.DayStart = 0;
          form.DayEnd = 0;
        }
      }}/>

                <label class='checkbox-label' for="AddStartOrEnd">Приключвам или започвам дейност този месец</label>
            </div>
        </div>

        <fieldset id="AddStartOrEndFieldset" style="display: none;">
            <div class="form-row">
                <div class="form-group">
                    <label for="DayStart">Начален ден на дейност</label>
                    <input class="input" id="DayStart" type="number" min="0" max="31" bind:value={form.DayStart}/>
                    <div class="info">Само ако започваш дейност през този месец</div>
                </div>

            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="DayEnd">Краен ден на дейност</label>
                    <input class="input" id="DayEnd" type="number" min="0" max="31" bind:value={form.DayEnd}/>
                    <div class="info">Само ако приключваш дейност през този месец</div>
                </div>
            </div>
        </fieldset>

        <div class="form-row">
            <div class="form-group submit-group">
                <button class="btn btn-large" on:click={saveIncome}>
                    <span><Save color="#444" size="20"/> Запази</span>
                </button>
            </div>
        </div>
        {#if data}
            <div class="alert success">
                <CircleCheck color="#748733" size="20"/>
                {data}
            </div>

            <div>
                <button class="btn" on:click={generateDeclarationOne}>
                    <span><BookText color="#444" size="20"/> Генерирай Декларация 1</span>
                </button>
            </div>

        {/if}
    </div>

    <div class="clearfix"></div>

   <Info />
</main>

<style>
    #home-page-block-right {
        float:right;
        padding-left:20px;
        width: 350px;
    }

    #home-page-input-box label {
        width: 160px;
    }

    #home-page-input-box div.info {
        padding-left: 170px;
    }

    #home-page-input-box .checkbox-group {
        padding-left: 170px;
    }

    #home-page-input-box label.checkbox-label {
        width: auto;
    }


    #tax-calculator-box, #declaration-six-box {
        padding-top: 20px;
    }

    .hidden-form-block {
        padding: 10px 0;
    }

    .hidden-form-block input[type="number"] {
        width: 100px;
    }

    #tax-calculator-result {
        margin: 10px 0;
    }

    #tax-enter-form-amount {
        width:70px;
    }
</style>
