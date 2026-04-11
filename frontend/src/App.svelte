<script lang="ts">
  import Router from "svelte-spa-router";
  import Home from './routes/Home.svelte';

  import Settings from './routes/Settings.svelte'
  import Personal from './routes/Personal.svelte'
  import Data from './routes/Data.svelte'

  import { House, User, Settings2, NotebookText } from 'lucide-svelte';
</script>


<div class="nav">
    <a href="#/"><House size="16" /> Начало</a>
    <a href="#/data"><NotebookText size="16" /> Въведени данни</a>
    <a href="#/personal"><User size="16" /> Лични данни</a>
    <a href="#/settings"><Settings2 size="16" /> Настройки</a>
</div>

<Router
    routes={{
        "/": Home,
        "/settings": Settings,
        "/personal": Personal,
        "/data": Data,
    }}
/>

<div class="footer">
  <p>Made with: Wails, Svelte, Lucide. <a href="https://github.com/sajax88/bumashtina" target="_blank">GitHub</a></p><!-- TODO: link-->
</div>