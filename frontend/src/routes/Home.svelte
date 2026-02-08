<script lang="ts">
  import { GenerateDeclarationOne } from "../../wailsjs/go/main/App.js";
  import { SaveConfig } from "../../wailsjs/go/main/App.js";
  import { LoadConfig } from "../../wailsjs/go/main/App.js";

  import { onMount } from 'svelte';

	onMount(() => {
    load_config()
	});

  let res = ""
  let config;

  function decl1(): void {
    GenerateDeclarationOne().then((result) => (res = result));
  }

  function save_config(): void {
    SaveConfig(config).then((result) => (res = result));
  }

  function load_config(): void {
    LoadConfig().then(function (result) {
        config = result
    });
  }
</script>

<main>
  <div class="input-box" id="input-box">
    {#if config}
    <input class="input" id="FirstName" bind:value={config.FirstName} type="text" />
    <input class="input" id="MiddleName" bind:value={config.MiddleName} type="text" />
    <input class="input" id="LastName" bind:value={config.LastName} type="text" />


     <button class="btn" on:click={save_config}>Save Config</button>
    {/if}

    <button class="btn" on:click={decl1}>Decl1</button>
    <div>{res}</div>

    <a href="#/page">Page</a>
  </div>
</main>

<style>
  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }
</style>
