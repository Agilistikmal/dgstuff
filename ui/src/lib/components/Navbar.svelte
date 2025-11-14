<script>
  import { Link } from "svelte-routing";
  import Logo from "./Logo.svelte";

</script>

<nav class="sticky top-0 z-50 w-full bg-white/5 backdrop-blur-sm">
  <div class="container mx-auto px-4 py-4">
    <div class="flex justify-between items-center">
      <Link to="/">
        <Logo />
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
