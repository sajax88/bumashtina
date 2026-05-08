<script lang="ts">
    import {GenerateDeclarationOne} from "../../wailsjs/go/main/App.js";
    import { onMount } from 'svelte';
    import {BookText} from "lucide-svelte";
    export let params;

    let month = parseInt(params.month);
    let year = parseInt(params.year);

    onMount(() => {

    });

    function decl1(): void {
        GenerateDeclarationOne(month, year).then(function (result) {
            console.log(result)
            // TODO: show success message, or error
        });
    }
</script>

<main>
    <div class="input-box" id="input-box">
    <h2>Въведени данни за {month}/{year}</h2>


        <button class="btn btn-large" id="declaration-one-button" on:click={() => decl1()}>
            <span><BookText color="#444" size="20"/> Генерирай Декларация 1</span>
        </button>
    </div>
</main>

<style>

</style>