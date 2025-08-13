describe('Login Test', () => {
    it('Successful', () => {
        const test_username = "e2e_tester";
        const test_password = "password1234";

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

        cy.contains('Login successful!', { timeout: 10000 });

    })

    it('Unsuccessful', () => {
        const test_username = "e2e_tester80";
        const test_password = "password1234";
        const test_email = "e2e_tester@gmail.com"

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

        cy.contains('Login failed: Invalid credentials', { timeout: 10000 });

    })
})
