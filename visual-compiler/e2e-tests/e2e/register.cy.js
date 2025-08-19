describe('Registration Test', () => {
    it('Successful', () => {
        const test_username = `e2e_tester_${Date.now()}`;
        const test_password = "password1234";
        const test_email = `e2e_tester_${Date.now()}@gmail.com`;

        cy.visit('http://localhost:5173/auth-page');

        cy.get('.tab-nav').should('be.visible');
        cy.get('.tab-nav button').contains("Register").should('be.visible');
        cy.contains('button', 'Login').should('have.class', 'active');
        cy.contains('button', 'Register').should('not.have.class', 'active');

        cy.wait(1000); 
        cy.get('.tab-nav button').contains("Register").click();
        cy.get('.tab-nav button').contains("Register").should('have.class', 'active');

        cy.get('body').should('contain', 'Username').and('contain', 'Email');

        cy.get('#regUsername', { timeout: 10000 }).should('be.visible');
        cy.get('#regUsername').type(test_username);
		cy.get('#regEmail').type(test_email);
		cy.get('#regPassword').type(test_password);
		cy.get('#regConfirmPassword').type(test_password);

        cy.get('.icon-submit-btn').click();

        cy.contains('Account created successfully', {timeout: 1000});

    })

    it('Unsuccessful', () => {
        const test_username = `e2e_tester_${Date.now()}`;
        const test_password = "pass";
        const test_email = `e2e_tester_${Date.now()}@gmail.com`;

        cy.visit('http://localhost:5173/auth-page');

        cy.get('.tab-nav').should('be.visible');
        cy.get('.tab-nav button').contains("Register").should('be.visible');
        cy.contains('button', 'Login').should('have.class', 'active');
        cy.contains('button', 'Register').should('not.have.class', 'active');

        cy.wait(1000); 
        cy.get('.tab-nav button').contains("Register").click();
        cy.get('.tab-nav button').contains("Register").should('have.class', 'active');

        cy.get('body').should('contain', 'Username').and('contain', 'Email');

        cy.get('#regUsername', { timeout: 10000 }).should('be.visible');
        cy.get('#regUsername').type(test_username);
		cy.get('#regEmail').type(test_email);
		cy.get('#regPassword').type(test_password);
		cy.get('#regConfirmPassword').type(test_password);

        cy.get('.icon-submit-btn').click();

        cy.contains('Registration error', {timeout: 1000});

    })
})