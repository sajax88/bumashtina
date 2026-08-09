<script lang="ts">
    import {Check, Scale} from "lucide-svelte";
    import {GetActiveMonthsNumber, DoYearlyAlignment} from "../../wailsjs/go/main/App";
    import {onMount} from "svelte";

    let alignmentResult = ''

    let yearlyAlignmentForm = {
        Year: new Date().getFullYear(),
        ActiveMonths: 12, // Recalculated when the year is changed
    }

    function fetchActiveMonthsNumber(year: number): void {
        GetActiveMonthsNumber(year).then((result: number) => (yearlyAlignmentForm.ActiveMonths = result));
    }

    function displayYearlyAlignmentForm(): void {
        DoYearlyAlignment(yearlyAlignmentForm.Year).then((result: string) => (alignmentResult = result))
    }

    $: fetchActiveMonthsNumber(yearlyAlignmentForm.Year);

    onMount(() => {
        fetchActiveMonthsNumber(yearlyAlignmentForm.Year)
    });
</script>

<div id="yearly-alignment-box">
    <div class="form-row">
        <div class="form-group">
            <button class="btn"
                    on:click={() => {document.getElementById('yearly-alignment-block').style.display = 'block';}}>
                <span><Scale color="#444" size="20"/> Годишно изравняване</span>
            </button>

            <div id="yearly-alignment-block" class="hidden-form-block" style="display: none;">
                <small>Година</small> <input type="number" id="yearly-alignment-year" class="year-input"
                                             bind:value={yearlyAlignmentForm.Year}/>

                <small style="margin-left: 5px;">Активни месеци: {yearlyAlignmentForm.ActiveMonths}</small>

                <button class="btn btn-small" on:click={displayYearlyAlignmentForm}>
                    <span><Check color="#444" size="20"/></span>
                </button>

                {#if alignmentResult}
                    <div id="yearly-alignment-result-block">
                        {alignmentResult}
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>

<style>
    #yearly-alignment-box {
        padding-top: 20px;
    }

    #yearly-alignment-block {
        padding-top: 10px;
    }
</style>