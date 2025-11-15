<script>
  import { onMount } from "svelte";
  import { AppInfoApi } from "../api/appinfo";
  import Logo from "./Logo.svelte";
  import { Link } from "svelte-routing";

  let appInfo = $state(null);
  onMount(async () => {
    appInfo = await AppInfoApi.get();
  });
</script>

<footer
  class="bg-brand bg-linear-to-b from-brand to-black/50 text-white pt-8 pb-24"
>
  <div class="container mx-auto px-4 py-4">
    <div class="flex justify-between items-center">
      <div>
        <div
          class="p-4 bg-white text-black rounded-lg flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <span class="font-medium">Powered by</span>
            <a href="https://agil.zip/app/dgstuff" target="_blank">
              <Logo />
            </a>
          </div>
          <span
            class="text-sm text-brand px-3 py-1 bg-brand/10 rounded-full font-medium"
          >
            {appInfo?.version}
          </span>
        </div>
        <p>dgstuff is a platform for selling and buying digital products.</p>
        <p class="text-sm mt-2">
          &copy; 2025 {appInfo?.name}. All rights reserved.
        </p>
      </div>
    </div>
  </div>
</footer>
