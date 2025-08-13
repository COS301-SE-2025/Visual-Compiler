describe('Analyser Test', ()=> {
    beforeEach('Login user',()=>{
        const test_username = "e2e_tester";
        const test_password = "password1234";

        cy.visit('http://localhost:5173/auth-page');

        cy.get('#loginUsername').type(test_username);
		cy.get('#loginPassword').type(test_password);

        cy.get('.icon-submit-btn').click();
    })

    it('Analyser phase', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        //source code
        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').click();
        //lexer
        cy.get('button').contains('Lexer').click();
        cy.get('.canvas-container').contains('Lexer').dblclick();
        cy.get('.automaton-btn').contains('Regular Expression').click();
        cy.get('.default-toggle-btn').click();
        cy.get('.submit-button').contains('Submit').click();
        cy.get('.generate-button').contains('Generate Tokens').click();
        cy.get('button').contains('Return to Canvas').click();
        //parser
        cy.get('button').should('contain', 'Parser');
        cy.get('button').contains('Parser').click();
        cy.get('.canvas-container').contains('Parser').dblclick();
        cy.get('.grammar-header').contains('Context-Free Grammar').click();
        cy.get('.default-toggle-btn').click();
        cy.get('.submit-button').contains('Submit Grammar').click();
        cy.get('.submit-button').contains('Generate Syntax Tree').click();
        cy.get('button').contains('Return to Canvas').click();

        //analyser
        cy.get('button').should('contain', 'Analyser');
        cy.get('button').contains('Analyser').click();

        cy.get('.canvas-container').should('contain','Analyser');
        cy.get('.canvas-container').contains('Analyser').dblclick();

         cy.get('.heading').should('contain','ANALYSING');
        cy.get('.grammar-header').contains('Context-Free Grammar').click();
        cy.get('.default-toggle-btn').click();

        cy.get('.submit-button').contains('Submit All Rules').click();
        cy.get('.generate-button').contains('Generate Symbol Table').click();

        cy.get('.artifact-viewer').should('contain','Syntax Tree');
        cy.get('.artifact-viewer').should('not.have.class','empty-state');

  
    })

});