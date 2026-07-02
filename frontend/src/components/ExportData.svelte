<script lang="ts">
    import {CircleCheck, FileDown} from 'lucide-svelte';
    import {ExportData} from "../../wailsjs/go/main/App";
    import {fade} from 'svelte/transition';

    let resExport
    async function exportData() {
        ExportData().then((result) => {resExport = result});
    }
</script>

<div id="export-data">
    <button class="btn btn-large" on:click={exportData}>
        <span><FileDown color="#444" size="20"/> Експортирай данните и настройките</span>
    </button>

    {#if resExport}
        <div class="alert success" in:fade={{duration:300}} out:fade={{duration:800}}>
            <CircleCheck color="#748733" size="20"/>  {resExport}
        </div>
        {/if}
</div>

<style>
    #export-data {
        margin: 10px 0 30px 0;
    }
</style>