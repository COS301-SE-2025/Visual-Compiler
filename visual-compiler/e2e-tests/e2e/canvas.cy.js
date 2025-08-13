describe('Canvas Test', ()=> {
    beforeEach('Login user',()=>{
        const test_username = "e2e_tester";
        const test_password = "password1234";

        cy.visit('http://localhost:5173/auth-page');

        cy.get('#loginUsername').type(test_username);
		cy.get('#loginPassword').type(test_password);

        cy.get('.icon-submit-btn').click();
    })

    it('Add Source Code Node', () => {
        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        cy.get('button').should('contain', 'Source Code');
        cy.get('button').contains('Source Code').click();

        cy.get('.canvas-container').should('contain','Source Code');
        cy.get('.canvas-container').contains('Source Code').dblclick();

        cy.get('.code-input-container').should('contain', 'Source Code Input');
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').contains('Confirm Code');
        cy.get('.confirm-btn').click();

        cy.contains('Code confirmed and saved');

    })

    it('Add Lexer Node', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').click();

        cy.get('button').should('contain', 'Lexer');
        cy.get('button').contains('Lexer').click();

        cy.get('.canvas-container').should('contain','Lexer');

    })

    it('Add Parser Node', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        cy.get('button').should('contain', 'Parser');
        cy.get('button').contains('Parser').click();

        cy.get('.canvas-container').should('contain','Parser');
        
    })

    it('Add Analyser Node', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        cy.get('button').should('contain', 'Analyser');
        cy.get('button').contains('Analyser').click();

        cy.get('.canvas-container').should('contain','Analyser');
        
    })

    it('Add Translator Node', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        cy.get('button').should('contain', 'Translator');
        cy.get('button').contains('Translator').click();

        cy.get('.canvas-container').should('contain','Translator');
        
    })
})