<script lang="ts">

  const logoUrl = '/half_stack_phoenix_only.png';
</script>

<header class="navbar">
  <img src={logoUrl} alt="Visual Compiler logo" class="logo" />
  <h1 class="title">Visual Compiler</h1>
</header>

<style>
  .navbar {
    position: relative;
    height: 3.5rem;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0,0,0,0.1);
    display: flex;
    justify-content: center; 
    align-items: center;
  }

  .logo {
    position: absolute;     
    left: 1rem;
    width: 2.5rem;
    height: auto;
  }

  .title {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
    color: #041a47;
  
  }
</style>
