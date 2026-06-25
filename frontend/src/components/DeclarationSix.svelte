<script lang="ts">
    import {BookText, Check, Save} from "lucide-svelte";
    import {GenerateDeclarationSix, PreviewDeclarationSix} from "../../wailsjs/go/main/App";

    // TODO: in one file?
    const MONEY_DIVIDER = 100;

    let declarationSixValues = {
        PensionPartOne: 0,
        PensionPartTwo: 0,
        HealthInsurance: 0,
    };
    let declarationSixResult = "";

    let declarationSixForm = {
        Year: new Date().getFullYear(),
    }

    function previewDeclarationSix(): void {
        PreviewDeclarationSix(declarationSixForm.Year).then(function (result) {
            if (result.PensionPartOneCents) {
                declarationSixValues.PensionPartOne = result.PensionPartOneCents / MONEY_DIVIDER
                declarationSixValues.PensionPartTwo = result.PensionPartTwoCents / MONEY_DIVIDER
                declarationSixValues.HealthInsurance = result.HealthInsuranceCents / MONEY_DIVIDER
            }
        });
    }

    function generateDeclarationSix(): void {
        let sums = {
            PensionPartOneCents: moneyToCents(declarationSixValues.PensionPartOne),
            PensionPartTwoCents: moneyToCents(declarationSixValues.PensionPartTwo),
            HealthInsuranceCents: moneyToCents(declarationSixValues.HealthInsurance),
        }
        GenerateDeclarationSix(declarationSixForm.Year, sums).then((result) => (declarationSixResult = result));
    }

    function moneyToCents(sum: string): number
    {
        return parseInt(Math.round(parseFloat(sum) * MONEY_DIVIDER))
    }
</script>

<div id="declaration-six-box">
    <div class="form-row">
        <div class="form-group">
            <button class="btn"
                    on:click={() => {document.getElementById('declaration-six-block').style.display = 'block';}}>
                <span><BookText color="#444" size="20"/> Генерирай Декларация 6</span>
            </button>
            <div id="declaration-six-block" class="hidden-form-block" style="display: none;">
                Година <input type="number" id="tax-calculator-year" class="year-input"
                              bind:value={declarationSixForm.Year}/>

                <button class="btn btn-small" on:click={previewDeclarationSix}>
                    <span><Check color="#444" size="20"/></span>
                </button>
            </div>

            {#if declarationSixResult}
                {declarationSixResult}
            {/if}

            {#if declarationSixValues.PensionPartOne}
                <div class="alert alert-info" id="declaration-six-result">
                    <label class="paid-insurance-label">ДОО: </label><input
                        class="paid-insurance-input" type="text"
                        bind:value={declarationSixValues.PensionPartOne}
                /><br/>
                    <label class="paid-insurance-label">ДЗПО: </label><input
                        class="paid-insurance-input" type="text"
                        bind:value={declarationSixValues.PensionPartTwo}
                /><br/>
                    <label class="paid-insurance-label">НЗОК: </label><input
                        class="paid-insurance-input" type="text"
                        bind:value={declarationSixValues.HealthInsurance}
                />
                    <br/>
                    <button class="btn btn-small" id="button-generate-declaration" on:click={generateDeclarationSix}>
                        <span><Save color="#444" size="20"/>Запази файл</span>
                    </button>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    #declaration-six-box {
        padding-top: 20px;
    }

    #declaration-six-block {
        padding-top: 10px;
    }

    .paid-insurance-input {
        width: 60px;
    }

    label.paid-insurance-label {
        width: auto;
    }

    #button-generate-declaration {
        margin-top: 10px;
        margin-left: 0;
    }
</style>