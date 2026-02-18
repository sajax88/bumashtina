<script lang="ts">
  import { GenerateDeclarationOne } from "../../wailsjs/go/main/App.js";
  import { LoadUserConfig, LoadSettingsConfig } from "../../wailsjs/go/main/App.js";

  import { onMount } from 'svelte';

	onMount(() => {
    load_configs()
	});

  let res = ""
  let configUser
  let configSettings

  function decl1(): void {
    GenerateDeclarationOne().then((result) => (res = result));
  }

  function load_configs(): void {
    LoadUserConfig().then(function (result) {
        configUser = result
    });
    LoadSettingsConfig().then(function (result) {
        configSettings = result
    });
  }
</script>

<main>
  <div class="input-box" id="input-box">
    <button class="btn" on:click={decl1}>Generate Declaration 1</button>
    <div>{res}</div>

    <a href="#/personal">Personal</a>
    <a href="#/settings">Settings</a>
  </div>
</main>

<style>
  
</style>
