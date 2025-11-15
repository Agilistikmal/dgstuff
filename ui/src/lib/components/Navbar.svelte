<script>
  import { Link } from "svelte-routing";
  import Logo from "./Logo.svelte";
  import { AppInfoApi } from "../api/appinfo";
  import { onMount } from "svelte";

  let appInfo = $state(null);
  onMount(async () => {
    appInfo = await AppInfoApi.get();
  });
</script>

<nav class="sticky top-0 z-50 w-full bg-white/5 backdrop-blur-sm">
  <div class="container mx-auto px-4 py-4">
    <div class="flex justify-between items-center">
      <Link to="/">
        {#if appInfo?.logo_url}
          <img
            src={appInfo?.logo_url}
            alt={appInfo?.name}
            class="h-10 text-2xl font-bold tracking-tighter text-black"
          />
        {:else}
          <Logo />
        {/if}
      </Link>
      <div class="flex items-center gap-x-8">
        <Link to="/">
          <span>Home</span>
        </Link>
        <Link to="/#stuffs">
          <span>Stuffs</span>
        </Link>
        <Link to="/#contact">
          <span>Contact</span>
        </Link>
      </div>
    </div>
  </div>
</nav>
