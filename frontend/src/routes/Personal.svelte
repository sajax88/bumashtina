<script lang="ts">
    import {LoadUserConfig, SaveUserConfig} from "../../wailsjs/go/main/App.js";
    import {onMount} from 'svelte';
    import {CircleCheck, Save} from 'lucide-svelte';
    import {fade} from 'svelte/transition';

    onMount(() => {
        loadConfig()
    });

    let res = ""
    let config;

    function saveConfig(): void {
        SaveUserConfig(config).then(
            function (result) {
                res = result
                setTimeout(() => {
                    res = ""
                }, 2000)
            },
        );
    }

    function loadConfig(): void {
        LoadUserConfig().then(function (result) {
            config = result
        });
    }
</script>

<main>
    <div class="input-box" id="input-box">

        <h2>Лични данни</h2>

        {#if config}
            <div class="form-row">
                <div class="form-group">
                    <label for="FirstName">Име</label>
                    <input class="input" required id="FirstName" bind:value={config.FirstName} type="text"/>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="MiddleName">Презиме</label>
                    <input class="input" id="MiddleName" bind:value={config.MiddleName} type="text"/>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="LastName">Фамилия</label>
                    <input class="input" required id="LastName" bind:value={config.LastName} type="text"/>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="Egn">ЕГН</label>
                    <input class="input" required type="text" id="Egn" bind:value={config.Egn}/>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="Bulstat">Булстат</label>
                    <input class="input" required type="text" id="Bulstat" bind:value={config.Bulstat}/>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="Email">Email</label>
                    <input class="input" required type="email" id="Email" bind:value={config.Email}/> <small>(нужен за Декларация 6)</small>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="Phone">Телефон</label>
                    <input class="input" required type="text" id="Phone" bind:value={config.Phone}/> <small>(нужен за Декларация 6)</small>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group submit-group">
                    <button class="btn btn-large" on:click={saveConfig}>
                        <span><Save color="#444" size="20"/> Запази</span>
                    </button>
                </div>
            </div>

            {#if res}
                <div class="alert success" in:fade={{duration:300}} out:fade={{duration:800}}>
                    <CircleCheck color="#748733" size="20"/> {res}
                </div>
            {/if}
        {/if}
    </div>
</main>

<style>

</style>
