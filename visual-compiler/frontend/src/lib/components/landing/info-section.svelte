<script>
  let active_tab = 'about'; 

  const team_members = [
    {
      name: 'Sashen Inder Gajai',
      role: 'Project Manager & Services Engineer',
      image_url: '/sashen.jpg',
    },
    {
      name: 'Shinn-Ru Hung',
      role: 'UI Engineer & Integration Engineer',
      image_url: '/shinn-ru.jpg',
    },
    {
      name: 'Yuthika Harripersad',
      role: 'Services Engineer & Testing Engineer',
      image_url: '/tia.jpg',
    },
    {
      name: 'Keanu Ellary',
      role: 'Integration Engineer & DevOps Engineer',
      image_url: '/keanu.jpg',
    },
    {
      name: 'Devan de Wet',
      role: 'UI Engineer & Data Engineer',
      image_url: '/devan.jpg',
    },
  ];
</script>

<section class="info_section">
  <div class="info_container">
    <div class="tab_nav">
      <button
        class:active={active_tab === 'about'}
        on:click={() => (active_tab = 'about')}>About The Project</button
      >
      <button
        class:active={active_tab === 'client'}
        on:click={() => (active_tab = 'client')}>About EPI-USE Labs</button
      >
      <button
        class:active={active_tab === 'team'}
        on:click={() => (active_tab = 'team')}>Meet The Team</button
      >
    </div>

    <div class="tab_content">
      {#if active_tab === 'about'}
        <div class="content_block">
          <h3>An Educational Tool for a Complex World</h3>
          <p>
            The Visual Compiler is an educational platform designed to demystify the complex world of compilers. The core goal is to provide a visual, block-based system where users can construct their own compilers phase by phase—from lexing and parsing to analysis and translation. By allowing users to see the visual artifacts like syntax trees and tokens at each stage, the platform makes abstract concepts tangible. It removes the need for deep technical skill, letting students experiment and learn in a high-level, intuitive environment.
          </p>
        </div>
      {:else if active_tab === 'client'}
        <div class="content_block">
          <h3>Value Through Innovation</h3>
          <p>
            EPI-USE Labs is a global software company that partners with clients to optimize the performance, management, and security of their SAP® and SAP SuccessFactors® systems. With over four decades of experience, they create innovative, well-engineered solutions to help clients navigate complex challenges. As part of 'Group Elephant', EPI-USE Labs is also committed to going 'Beyond Corporate Purpose' by contributing 1% of its revenue to the preservation of at-risk elephants and rhinos through community upliftment.
          </p>
        </div>
      {:else if active_tab === 'team'}
        <div class="content_block team_grid">
          {#each team_members as member}
            <div class="team_card">
              <img
                src={member.image_url}
                alt="Photo of {member.name}"
                class="team_photo"
              />
              <h4 class="member_name">{member.name}</h4>
              <p class="member_role">{member.role}</p>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</section>

<style>
  .info_section {
    background-color: #fafafa;
    color: #041a47;
    padding: 6rem 2rem;
  }

  .info_container {
    max-width: 1200px; 
    margin: 0 auto;
  }

  .tab_nav {
    display: flex;
    justify-content: center;
    gap: 1rem;
    margin-bottom: 3rem;
  }

  .tab_nav button {
    font-family: 'Inter', sans-serif;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 600;
    border: 2px solid #041a47;
    background-color: transparent;
    color: #041a47;
    border-radius: 50px;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .tab_nav button:hover {
    background-color: rgba(4, 26, 71, 0.1);
  }

  .tab_nav button.active {
    background-color: #041a47;
    color: #fafafa;
  }

  .tab_content {
    background-color: #ffffff;
    padding: 3rem;
    border-radius: 16px;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.07);
    min-height: 300px;
  }
  
  .content_block h3 {
      font-size: 1.75rem;
      margin-bottom: 1rem;
      text-align: center;
  }

  .content_block p {
      font-size: 1.1rem;
      line-height: 1.7;
      text-align: left;
      margin-left: 1rem;
      margin-right: 1rem;
  }

  .team_grid {
      display: grid;
      grid-template-columns: repeat(5, 1fr); 
      gap: 1.5rem;
      align-items: start;
  }

  .team_card {
      display: flex;
      flex-direction: column;
      align-items: center;
  }

  .team_photo {
      width: 150px;
      height: 150px;
      border-radius: 50%;
      object-fit: cover;
      margin-bottom: 1rem;
      border: 4px solid #fafafa;
      box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  }
  

  .member_name {
      font-size: 1.2rem;
      font-weight: 600;
      margin-bottom: 0.5rem; 
      text-align: center;
  }

  .member_role {
      font-size: 0.9rem;
      opacity: 0.8;
      text-align: center;
      max-width: 180px; 
  }

  @media (max-width: 1200px) {
    .team_grid {
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 2rem;
    }
  }
</style>
