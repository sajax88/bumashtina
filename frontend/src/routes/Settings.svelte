<script lang="ts">
    import {
        LoadSettingsConfig,
        LoadTaxesConfig,
        LoadTaxesConfigLabels,
        SaveSettingsConfig, SaveTaxesConfig,
    } from "../../wailsjs/go/main/App.js";

    import {onMount} from 'svelte';

    import {CircleCheck, Save} from 'lucide-svelte';

    onMount(() => {
        load_config()
    });

    const MONEY_DIVIDER = 100;
    let res = ""
    let resTaxes = ""
    let settingsConfig;
    let taxesConfig;
    let taxesLabels;

    function saveSettingsConfig(): void {
        SaveSettingsConfig(settingsConfig).then((result) => (res = result));
    }

    function saveTaxesConfig(): void {
        SaveTaxesConfig(taxesConfig).then((result) => (resTaxes = result));
    }

    function load_config(): void {
        LoadSettingsConfig().then(function (result) {
            settingsConfig = result
        });
        LoadTaxesConfig().then(function (result) {
            taxesConfig = result
        });
        LoadTaxesConfigLabels().then(function (result) {
            taxesLabels = result
        });
    }
</script>

<main>
    <div class="input-box" id="input-box-settings">
        <h2>Настройки</h2>

        {#if settingsConfig}
            <div class="form-row">
                <div class="form-group checkbox-group">
                    <input class="checkbox" id="IsPregnancyInsuranceEnabled"
                           bind:checked={settingsConfig.IsPregnancyInsuranceEnabled} type="checkbox"/>
                    <label class="checkbox-label" for="IsPregnancyInsuranceEnabled">Общо заболяване и майчинство</label>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group submit-group">
                    <button class="btn btn-large" on:click={saveSettingsConfig}>
                        <span><Save color="#444" size="20"/> Запази</span>
                    </button>
                </div>
            </div>

            {#if res}
                <div class="alert success">
                    <CircleCheck color="#748733" size="20"/>  {res}
                </div>
            {/if}

        {/if}

        <h2>Данъци и осигуровки, 2026</h2>
        {#if taxesConfig && taxesLabels}
            <table class="table">
                <tbody>
                {#each Object.entries(taxesConfig) as [key, value]}
                    <tr>
                        <td>{taxesLabels[Object.keys(taxesConfig).indexOf(key)]}</td>
                        <td>
                            {#if Object.keys(taxesConfig).indexOf(key) < 2}
                                {value / MONEY_DIVIDER}
                            {:else}
                                {value} %
                            {/if}

                            <!-- TODO: edit button:
                            <input class="input taxes-config-input" id={key} type="text"
                                   bind:value={taxesConfigInputValues[key]}
                                   on:input={() => updateTaxesConfigValue(key, index)}/>
                            {#if !isMoneyField(index)}
                                <span class="taxes-config-unit">%</span>
                            {/if}
                            -->
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>

            <div class="submit-group">
                <button class="btn btn-large" on:click={saveTaxesConfig}>
                    <span><Save color="#444" size="20"/> Запази</span>
                </button>
            </div>

            {#if resTaxes}
                <div class="alert success">
                    <CircleCheck color="#748733" size="20"/> {resTaxes}
                </div>
            {/if}
        {/if}
    </div>
</main>

<style>
    .table td:first-child {
        width: 350px;
    }

    .table td:last-child {
        font-weight: bold;
    }
</style>