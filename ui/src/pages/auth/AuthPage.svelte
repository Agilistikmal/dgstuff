<script>
  import { onMount } from "svelte";
  import { AppInfoApi } from "../../lib/api/appinfo";
  import Footer from "../../lib/components/Footer.svelte";
  import Loading from "../../lib/components/Loading.svelte";
  import Navbar from "../../lib/components/Navbar.svelte";
  import LoginForm from "./LoginForm.svelte";
  import RegisterForm from "./RegisterForm.svelte";

  let queryParams = $state(new URLSearchParams(window.location.search));
  let tab = $state(queryParams.get("tab") || "login");
</script>

<Navbar />
<div class="min-h-screen w-full max-w-screen-sm mx-auto p-4 pb-24">
  {#await AppInfoApi.get()}
    <Loading />
  {:then appInfo}
    {#if appInfo.is_first_launch || tab === "register"}
      <RegisterForm />
    {:else}
      <LoginForm />
    {/if}
  {:catch error}
    <div>
      <span class="text-red-500">{error.message}</span>
    </div>
  {/await}
</div>
<Footer />
