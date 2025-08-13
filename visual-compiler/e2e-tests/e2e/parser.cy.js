describe('Parser Test', ()=> {
    beforeEach('Login user',()=>{
        const test_username = "e2e_tester";
        const test_password = "password1234";

        cy.visit('http://localhost:5173/auth-page');

        cy.get('#loginUsername').type(test_username);
		cy.get('#loginPassword').type(test_password);

        cy.get('.icon-submit-btn').click();
    })

    it('Parser phase', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').click();

        cy.get('button').contains('Lexer').click();
        cy.get('.canvas-container').contains('Lexer').dblclick();
        cy.get('.automaton-btn').contains('Regular Expression').click();
        cy.get('.default-toggle-btn').click();
        cy.get('.submit-button').contains('Submit').click();
        cy.get('.generate-button').contains('Generate Tokens').click();
        cy.get('button').contains('Return to Canvas').click();

        cy.get('button').should('contain', 'Parser');
        cy.get('button').contains('Parser').click();
        cy.get('.canvas-container').contains('Parser').dblclick();

        cy.get('.phase-inspector').should('contain','PARSING');
        cy.get('.grammar-header').contains('Context-Free Grammar').click();
        cy.get('.default-toggle-btn').click();

        cy.get('.submit-button').contains('Submit Grammar').click();
        cy.get('.submit-button').contains('Generate Syntax Tree').click();

        cy.get('.artifact-viewer').should('contain','Syntax Tree');
        cy.get('.artifact-viewer').should('not.have.class','empty-state');
  
    })

});