<script lang="ts">
    import {CircleCheck, FileUp} from "lucide-svelte";
    import {fade} from 'svelte/transition';
    import {
        ImportData
    } from "../../wailsjs/go/main/App.js";

    let resImport;
    async function importData() {
        ImportData().then((result) => {resImport = result});
    }
</script>

<div>
    <button class="btn btn-large" on:click={importData}>
        <span><FileUp color="#444" size="20"/> Импортирай данни и настройки</span>
    </button>


    {#if resImport}
        <div class="alert success" in:fade={{duration:300}} out:fade={{duration:800}}>
            <CircleCheck color="#748733" size="20"/>  {resImport}
        </div>

    {:else}
        <div class="alert warning"><small>
            Внимание, при импорт на данните ще се заличат въведените ви данни и настройки!
        </small></div>

    {/if}
</div>