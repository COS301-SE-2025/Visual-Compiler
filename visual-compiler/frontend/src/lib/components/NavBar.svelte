<script lang="ts">
  // If you placed the PNG in your `static/` folder:
  const logoUrl = '/half_stack_phoenix_only.png';
</script>

<header class="navbar">
  <div class="logo-container">
    <img src={logoUrl} alt="Visual Compiler logo" class="logo" />
    <h1>Visual Compiler</h1>
  </div>
</header>

<style>
  .navbar {
    height: 3.5rem;
    background: #fff;
    display: flex;
    align-items: center;
    padding: 0 1rem;
    box-shadow: 0 1px 4px rgba(0,0,0,0.1);
  }

  .logo-container {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .logo {
    width: 2.5rem;          
    height: auto;
    margin-right:1rem;
  }

  .navbar h1 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
    color: #041a47;
  }
</style>
