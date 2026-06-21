<script lang="ts">
    import {
        LoadSettingsConfig,
        LoadTaxesConfig,
        LoadTaxesConfigLabels,
        SaveSettingsConfig, SaveTaxesConfig,
    } from "../../wailsjs/go/main/App.js";

    import {onMount} from 'svelte';

    import {CircleCheck, Edit, Save} from 'lucide-svelte';

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
        let taxesConfigToSave = {}
        for(const k in taxesConfig) {
            if (k > 1) {
                // Percentage
                taxesConfigToSave[k] = parseInt(taxesConfig[k])
            } else {
                taxesConfigToSave[k] = parseFloat(taxesConfig[k])
            }
        }
        SaveTaxesConfig(taxesConfigToSave).then((result) => (resTaxes = result));
        // TODO: hide inputs!
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

    function showEditInput(key) {
        // TODO: money divider for k < 2
        document.getElementById(`taxes-config-input-${key}`).style.display = "block";
        document.getElementById(`taxes-config-display-value-${key}`).style.display = "none";
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
                    <p><small><i>
                        Самоосигуряващи се лица избират дали да се осигуряват за общо заболяване и майчинство
                        при подаване на "Декларация за регистрация на самоосигуряващо се лице". Преди да промените тази настройка,
                        убедете се, че сте подали такава Декларация в НАП, тъй като настройката ще промени данните
                        в следващите ви Декларации 1.
                    </i></small></p>
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
                            <span class="taxes-config-display-value" id="taxes-config-display-value-{key}">
                            {#if Object.keys(taxesConfig).indexOf(key) < 2}
                                {value / MONEY_DIVIDER}
                            {:else}
                                {value} %
                            {/if}
                            </span>

                            <input style="display: none" id="taxes-config-input-{key}" class="input taxes-config-input" type="text"
                                   bind:value={taxesConfig[key]}
                            />
                        </td>
                        <td>
                            <button class="btn btn-small" on:click="{() => showEditInput(key)}">
                                <Edit color="#444" size="20"/>
                            </button>
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