<script lang="ts">
    import {Calculator, Check, Plus} from "lucide-svelte";

    import {
        CalculateTaxForQuarter,
        SavePaidTaxForQuarter
    } from "../../wailsjs/go/main/App.js";

    const MONEY_DIVIDER = 100;

    let taxCalculationResult;
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

                <input type="number" id="tax-calculator-year" class="year-input" bind:value={taxCalculatorForm.Year}/>

                <button class="btn btn-small" on:click={() => {calculateTaxForQuarter()}}>
                    <span><Check color="#444" size="20"/></span>
                </button>
                {#if taxCalculationResult}
                    <div class="alert alert-info" id="tax-calculator-result">
                        <p>Месеци: {taxCalculationResult.MonthStart}-{taxCalculationResult.MonthEnd}</p>
                        <p>({taxCalculationResult.TotalIncomeCents / MONEY_DIVIDER} доход - {taxCalculationResult.ExpensesCents / MONEY_DIVIDER}
                        приспадащи се разходи - {taxCalculationResult.PaidInsuranceCents / MONEY_DIVIDER} платени осигуровки) *
                        10% =
                            <b>{taxCalculationResult.TaxCents / MONEY_DIVIDER} EUR</b></p>
                        <!-- TODO: take tax % from dynamic -->
                        <!-- TODO: поясни, откуда сумма расходов - 25% !-->
                    </div>
                    <!-- TODO: warn if not entered
                    <small><i>Точният данък зависи от фактически платени осигуровки</i></small>
                    -->
                {/if}
            </div>
        </div>
    </div>
</div>

<div id="tax-enter-form-box">
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

                <input type="number" id="tax-enter-form-year" class="year-input" bind:value={taxEnterForm.Year}/>

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

<style>
    #tax-calculator-box, #tax-enter-form-box {
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