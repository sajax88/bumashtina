<script lang="ts">
    import {BookText, Check} from "lucide-svelte";
    import {GenerateDeclarationSix} from "../../wailsjs/go/main/App";

    let declarationSixResult = ""; // TODO: if err

    let declarationSixForm = {
        Year: new Date().getFullYear(),
    }

    function generateDeclarationSix(): void {
        GenerateDeclarationSix(declarationSixForm.Year).then((result) => (declarationSixResult = result));
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
                Година <input type="number" id="tax-calculator-year" class="year-input" bind:value={declarationSixForm.Year}/>

                <button class="btn btn-small" on:click={generateDeclarationSix}>
                    <span><Check color="#444" size="20"/></span>
                </button>
            </div>
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
</style>