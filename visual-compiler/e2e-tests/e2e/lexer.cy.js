describe('Lexer Test', ()=> {
    beforeEach('Login user',()=>{
        const test_username = "e2e_tester";
        const test_password = "password1234";

        cy.visit('http://localhost:5173/auth-page');

        cy.get('#loginUsername').type(test_username);
		cy.get('#loginPassword').type(test_password);

        cy.get('.icon-submit-btn').click();
    })

    it('Lexer phase', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').click();

        cy.get('button').contains('Lexer').click();
        cy.get('.canvas-container').contains('Lexer').dblclick();

        cy.get('.phase-inspector').should('contain','LEXING');
        cy.get('.automaton-btn').contains('Regular Expression').click();
        cy.get('.default-toggle-btn').click();

        cy.get('.submit-button').contains('Submit').click();
        cy.get('.generate-button').contains('Generate Tokens').click();

        cy.get('.token-table').should('contain','Type');
        cy.get('.token-table').should('contain','Value');
        cy.get('.token-table').should('contain','KEYWORD');
        cy.get('.token-table').should('contain','int');
        cy.get('.token-table').should('contain','IDENTIFIER');
        cy.get('.token-table').should('contain','blue');
        cy.get('.token-table').should('contain','ASSIGNMENT');
        cy.get('.token-table').should('contain','=');
        cy.get('.token-table').should('contain','INTEGER');
        cy.get('.token-table').should('contain','13');
        cy.get('.token-table').should('contain','OPERATOR');
        cy.get('.token-table').should('contain','22');
        cy.get('.token-table').should('contain','SEPARATOR');
        cy.get('.token-table').should('contain',';');
  
    })

});