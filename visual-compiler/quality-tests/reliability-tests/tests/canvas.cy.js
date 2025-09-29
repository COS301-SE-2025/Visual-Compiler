describe('Canvas Test', ()=> {
   beforeEach('Login user',()=>{
        const test_username = "e2e tester";
        const test_password = "testUser13";

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

        cy.contains('Welcome');
    })

    it('No API connection', () => {

        cy.wait(8000);
        cy.visit('http://localhost:5173/main-workspace');

        //delete project
        cy.get('.section-heading').should('contain', 'Start a new project');
        cy.contains('Service temporarily unavailable. Please try again later.');
        cy.wait(500);

    })

})
