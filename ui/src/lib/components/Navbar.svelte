<script>
  import Logo from "./Logo.svelte";

</script>

<nav class="sticky top-0 z-50 w-full bg-white/5 backdrop-blur-sm">
  <div class="container mx-auto px-4 py-4">
    <div class="flex justify-between items-center">
      <a href="/">
        <Logo />
      </a>
      <div class="flex items-center gap-x-8">
        <a href="/">
          <span>Home</span>
        </a>
        <a href="/#stuffs">
          <span>Stuffs</span>
        </a>
        <a href="/#contact">
          <span>Contact</span>
        </a>
      </div>
    </div>
  </div>
</nav>
