<script>
	let active_tab = 'about';

	const team_members = [
		{
			name: 'Sashen Inder Gajai',
			role1: 'Project Manager',
			role2: 'Services Engineer',
			image_url: '/sashen.jpg'
		},
		{
			name: 'Shinn-Ru Hung',
			role1: 'UI Engineer',
			role2: 'Integration Engineer',
			image_url: '/austin.jpg'
		},
		{
			name: 'Tia Harripersad',
			role1: 'Services Engineer',
			role2: 'Testing Engineer',
			image_url: '/tia.jpg'
		},
		{
			name: 'Keanu Ellary',
			role1: 'Integration Engineer',
			role2: 'DevOps Engineer',
			image_url: '/keanu.jpg'
		},
		{
			name: 'Devan de Wet',
			role1: 'UI Engineer',
			role2: 'Data Engineer',
			image_url: '/devan.jpg'
		}
	];
</script>

<section class="info_section">
	<div class="info_container">
		<div class="tab_nav">
			<button class:active={active_tab === 'about'} on:click={() => (active_tab = 'about')}
				>About The Project</button
			>
			<button class:active={active_tab === 'team'} on:click={() => (active_tab = 'team')}
				>Meet The Team</button
			>
			<button class:active={active_tab === 'client'} on:click={() => (active_tab = 'client')}
				>About EPI-USE Labs</button
			>
		</div>

		<div class="tab_content">
			{#if active_tab === 'about'}
				<div class="content_block">
					<h3>Visual Compiler</h3>
					<div class="new_info">
						<p>
							Visual Compiler is an educational platform that demonstrates the four main phases of compilation, specifically lexical analysis, syntax analysis, semantic analysis and code generation.
						</p>
						<p>
							Through integrated simulators, visualisers and tutorials, this web application fosters a constructive environment for computer science students to bridge theoretical concepts and practical application.
						</p>
					</div>
				</div>
			{:else if active_tab === 'client'}
				<div class="content_block">
					<h3>Value Through Innovation</h3>
					<div class="new_info">
						<p>
							EPI-USE Labs is a global software company that partners with clients to create innovative, well-engineered solutions that optimise the performance, management, and security of their SAP® and SAP SuccessFactors® systems.
						</p>
						<p>
							As part of Group Elephant, EPI-USE Labs is also committed to going beyond corporate purpose by contributing 1% of its revenue to the preservation of elephants and rhinos through community upliftment.
						</p>
					</div>
				</div>
			{:else if active_tab === 'team'}
				<div class="content_block team_grid">
					{#each team_members as member}
						<div class="team_card">
							<img src={member.image_url} alt="Photo of {member.name}" class="team_photo" />
							<h4 class="member_name">{member.name}</h4>
							<p class="member_role">{member.role1}</p>
							<p class="member_role">{member.role2}</p>
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

	.new_info {
		margin-top: 2rem;
		margin-bottom: 1rem;
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
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
	}

	.member_name {
		font-size: 1.2rem;
		font-weight: 600;
		margin-bottom: 0.5rem;
		text-align: center;
	}

	.member_role {
		margin:0;
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
