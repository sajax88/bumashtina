<script lang="ts">
    import {
        LoadSettingsConfig,
        LoadTaxesConfig,
        LoadTaxesConfigLabels,
        SaveSettingsConfig, SaveTaxesConfig,
    } from "../../wailsjs/go/main/App.js";
    import {onMount} from 'svelte';
    import {CircleCheck, Edit, Save} from 'lucide-svelte';
    import {fade} from 'svelte/transition';
    import {MONEY_DIVIDER} from '../constants'

    onMount(() => {
        loadAllConfigs()
    });

    // TODO: from some common source?

    let resSettings = ""
    let resTaxes = ""
    let settingsConfig;
    let taxesConfig;
    let taxesLabels;

    function saveSettingsConfig(): void {
        SaveSettingsConfig(settingsConfig).then(
            function (result) {
                resSettings = result
                setTimeout(() => {
                    resSettings = ""
                }, 2000)
            },
        );
    }

    function saveTaxesConfig(): void {
        let taxesConfigToSave = {}
        for(const k in taxesConfig) {
            if (k == 'MinInsuranceIncomeCents' || k == 'MaxInsuranceIncomeCents') {
                // Percentage
                taxesConfigToSave[k] = parseInt(Math.round(parseFloat(taxesConfig[k]) * MONEY_DIVIDER))
            } else {
                taxesConfigToSave[k] = parseFloat(taxesConfig[k])
            }
        }
        SaveTaxesConfig(taxesConfigToSave).then(
            function (result) {
                if (result == "") {
                    return
                }

                // Hide all inputs again
                const inputsList = document.querySelectorAll<HTMLElement>(".taxes-config-input");
                for (let i = 0; i < inputsList.length; i++) {
                    inputsList[i].style.display = "none";
                }

                const valuesList = document.querySelectorAll<HTMLElement>(".taxes-config-display-value");
                for (let i = 0; i < valuesList.length; i++) {
                    valuesList[i].style.display = "block";
                }

                resTaxes = result
                setTimeout(() => {
                    resTaxes = ""
                }, 2000)
            },
        );
    }

    function loadAllConfigs(): void {
        LoadSettingsConfig().then(function (result) {
            settingsConfig = result
        });
        LoadTaxesConfig().then(function (result) {
            taxesConfig = result
            taxesConfig.MinInsuranceIncomeCents = parseFloat(taxesConfig.MinInsuranceIncomeCents) / MONEY_DIVIDER
            taxesConfig.MaxInsuranceIncomeCents = parseFloat(taxesConfig.MaxInsuranceIncomeCents) / MONEY_DIVIDER
        });
        LoadTaxesConfigLabels().then(function (result) {
            taxesLabels = result
        });
    }

    function toggleEditInput(key: string) {
        if (document.getElementById(`taxes-config-input-${key}`).style.display == "block") {
            document.getElementById(`taxes-config-input-${key}`).style.display = "none";
            document.getElementById(`taxes-config-display-value-${key}`).style.display = "block";
        } else{
            document.getElementById(`taxes-config-input-${key}`).style.display = "block";
            document.getElementById(`taxes-config-display-value-${key}`).style.display = "none";
        }
    }
</script>

<main>
    <div class="input-box" id="input-box-settings">
        <h2>Настройки</h2>

        {#if settingsConfig}
            <div class="form-row">
                <p><small><i>
                    Самоосигуряващи се лица избират дали да се осигуряват за общо заболяване и майчинство
                    при подаване на "Декларация за регистрация на самоосигуряващо се лице". Преди да промените тази настройка,
                    убедете се, че сте подали такава Декларация в НАП, тъй като настройката ще промени данните
                    в следващите ви Декларации 1.
                </i></small></p>
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

            {#if resSettings}
                <div class="alert success" in:fade={{duration:300}} out:fade={{duration:800}}>
                    <CircleCheck color="#748733" size="20"/>  {resSettings}
                </div>
            {/if}
        {/if}

        <h2>Данъци и осигуровки</h2>
        <p><small><i>
            Промяна на тези настройки ще се отрази на новите ви Декларации 1, но няма да засегне вече въведените данни.
        </i></small></p>
        {#if taxesConfig && taxesLabels}
            <table class="table">
                <tbody>
                {#each Object.entries(taxesConfig) as [key, value]}
                    <tr>
                        <td>{taxesLabels[Object.keys(taxesConfig).indexOf(key)]}</td>
                        <td>
                            <span class="taxes-config-display-value" id="taxes-config-display-value-{key}">
                            {#if Object.keys(taxesConfig).indexOf(key) < 2}
                                {value} EUR
                            {:else}
                                {value} %
                            {/if}
                            </span>

                            <input style="display: none" id="taxes-config-input-{key}" class="input taxes-config-input" type="text"
                                   bind:value={taxesConfig[key]}
                            />
                        </td>
                        <td>
                            <button class="btn btn-small" on:click="{() => toggleEditInput(key)}">
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
                <div class="alert success" in:fade={{duration:300}} out:fade={{duration:800}}>
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

    .taxes-config-input {
        width: 80px;
    }
</style>