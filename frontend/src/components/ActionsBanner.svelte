<script lang="ts">
    import {CloudAlert} from 'lucide-svelte';
    import {onMount} from "svelte";
    import {LoadThisMonthActions} from "../../wailsjs/go/main/App";

    let thisMonthActions;

    onMount(async () => {
        await loadActions()
    });

    async function loadActions(): Promise<void> {
        const [actions] = await Promise.all([
            LoadThisMonthActions()
        ]);
        thisMonthActions = actions;
    }
</script>

{#if thisMonthActions}
    <div id="this-month-actions">
        <div class="alert success">
            <CloudAlert color="#778a3b" size="28"/> {thisMonthActions}
        </div>
    </div>
{/if}

<style>
    #this-month-actions {
        padding: 0 30px;
    }
</style>