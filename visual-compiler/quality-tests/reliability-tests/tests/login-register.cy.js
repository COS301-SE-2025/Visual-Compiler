describe('No API Connnection', () => {
    it('Unsuccessful Register', () => {
        const test_username = `e2e_tester_${Date.now()}`;
        const test_password = "pass";
        const test_email = `e2e_tester_${Date.now()}@gmail.com`;

        cy.visit('http://localhost:5173/auth-page');

        cy.get('.tab-nav').should('be.visible');
        cy.get('.tab-nav button').contains("Register").should('be.visible');
        cy.contains('button', 'Login').should('have.class', 'active');
        cy.contains('button', 'Register').should('not.have.class', 'active');

        cy.wait(500); 
        cy.get('.tab-nav button').contains("Register").click();
        cy.get('.tab-nav button').contains("Register").should('have.class', 'active');

        cy.get('body').should('contain', 'Username').and('contain', 'Email');

        cy.get('#regUsername', { timeout: 10000 }).should('be.visible');
        cy.get('#regUsername').type(test_username);
		cy.get('#regEmail').type(test_email);
		cy.get('#regPassword').type(test_password);
		cy.get('#regConfirmPassword').type(test_password);

        cy.get('.icon-submit-btn').click();

        cy.contains('Service temporarily unavailable. Please try again later.');

    })

    it('Unsucessful login', () => {
        const test_username = "e2e teste";
        const test_password = "testUser1";

        cy.visit('http://localhost:5173/auth-page');

        cy.get('.tab-nav').should('be.visible');
        cy.get('.tab-nav button').contains("Register").should('be.visible');
        cy.contains('button', 'Login').should('have.class', 'active');
        cy.contains('button', 'Register').should('not.have.class', 'active');

        cy.get('body').should('contain', 'Username');

        cy.get('#loginUsername', { timeout: 10000 }).should('be.visible');
        cy.get('#loginUsername').type(test_username);
		cy.get('#loginPassword').type(test_password);

        cy.get('.icon-submit-btn').click();

        cy.contains('Service temporarily unavailable. Please try again later.');

    })
})

