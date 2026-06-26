<script lang="ts">
    import {CloudAlert} from 'lucide-svelte';
    import {onMount} from "svelte";
    import {LoadAlerts} from "../../wailsjs/go/main/App";

    let alerts:String;

    onMount(async () => {
        await loadAlerts()
    });

    async function loadAlerts(): Promise<void> {
        LoadAlerts().then((result) => {alerts = result});
    }
</script>

{#if alerts}
    <div id="top-alerts">
        <div class="alert warning">
            <CloudAlert class="icon-success" color="#778a3b" size="28" style="display: none"/>
            <CloudAlert class="icon-warning" color="#b87a07" size="28" />
            <span>{alerts}</span>
        </div>
    </div>
{/if}

<style>
    #top-alerts {
        padding: 0 20px;
    }

    #top-alerts .alert {
        display: flex;
        align-items: center;
        color: #6e4905;
    }
</style>