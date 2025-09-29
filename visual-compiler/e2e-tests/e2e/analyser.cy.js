/*describe('Analyser Test', ()=> {
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

    it('Analyser phase', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        const project_name = 'canvas_project';
        //delete project
        cy.get('.section-heading').should('contain', 'Start a new project');
        cy.get('.project-block').contains(project_name).get('.delete-button').click();
        cy.get('.delete-confirm-button').click();
        cy.wait(500);

        //select project
        cy.get('.section-heading').should('contain', 'Start a new project');
        cy.get('.project-button').contains('New Blank').click();
        cy.wait(500);
        cy.get('#project-name-input').type(project_name);
        cy.get('#confirm-project-name').click();
        cy.wait(500);

        cy.get('button').contains("Next").click();
        cy.get('button').contains("Next").click();
        cy.get('button').contains("Next").click();
        cy.get('button').contains("Next").click();
        cy.get('button').contains("Next").click();
		cy.get('button').contains("Next").click();
        cy.get('button').contains("Get Started").click();

        //source code node
        cy.get('button').should('contain', 'Source Code');
        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').should('contain','Source Code');
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.wait(500);
        cy.get('.code-input-container').should('contain', 'Source Code Input');
        cy.get('button').contains('Show Example').click();
        cy.get('.confirm-btn').contains('Confirm Code');
        cy.get('.confirm-btn').click();
        cy.contains('Source code saved');

         //lexer node
        cy.get('button').should('contain', 'Lexer');
        cy.get('button').contains('Lexer').click();
        cy.get('.canvas-container').should('contain','Lexer');

        //parser node
        cy.wait(500);
        cy.get('button').should('contain', 'Parser');
        cy.get('button').contains('Parser').click();
       
        cy.get('button').should('contain', 'Analyser');
        cy.get('button').contains('Analyser').click();
        cy.get('.canvas-container').should('contain','Analyser');

        cy.wait(500);
        cy.get('#A-1\\/N-source-1').trigger('mousedown', { which: 1, force: true });
        cy.get('#A-1\\/N-lexer-2').trigger('mousemove', { force: true }).trigger('mouseup', { force: true });
        cy.wait(500);
        cy.get('#A-2\\/N-lexer-2').trigger('mousedown', {which: 1, force: true});
        cy.get('#A-1\\/N-parser-3').trigger('mousemove', {force: true}).trigger('mouseup', {force: true});
        cy.wait(500);
        cy.get('#A-2\\/N-parser-3').trigger('mousedown', {which: 1, force: true});
        cy.get('#A-1\\/N-analyser-4').trigger('mousemove', {force: true}).trigger('mouseup', {force: true});
        cy.wait(500);

        cy.get('.canvas-container').contains('Lexer').dblclick();
        cy.wait(500);
        cy.get('.phase-inspector').should('contain','LEXING');
        cy.get('.automaton-btn').contains('Regular Expression').click();
        cy.get('button').contains('Show Example').click();
        cy.get('.submit-button').contains('Submit').click();
        cy.get('.generate-button').contains('Generate Tokens').click();
        cy.get('.return-button').click();

        cy.wait(1000);
        cy.get('.canvas-container').contains('Parser').dblclick();
        cy.wait(500);
        cy.get('.phase-inspector').should('contain','PARSING');
        cy.get('.grammar-header').should('contain','Context-Free Grammar');
        cy.get('button').contains('Show Example').click();
        cy.get('.submit-button').contains('Submit Grammar').click();
        cy.get('.submit-button').contains('Generate Syntax Tree').click();
        cy.get('.return-button').click();

        //analyser
        cy.wait(1000);
        cy.get('.canvas-container').contains('Analyser').dblclick();
        cy.wait(500);
        cy.get('.canvas-container').contains('Analyser').dblclick();

        cy.get('.analyser-heading').should('contain','ANALYSING');
        cy.get('button').contains('Show Example').click();

        cy.get('.submit-button').contains('Submit Rules').click();
        cy.get('.generate-button').contains('Generate Symbol Table').click();

        cy.get('.artifact-viewer').should('contain','Symbol Table');
        cy.get('.symbol-table').should('contain','Type');
        cy.get('.symbol-table').should('contain','Name');
        cy.get('.symbol-table').should('contain','Scope');
        cy.get('.symbol-table').should('contain','int');
        cy.get('.symbol-table').should('contain','blue');
        cy.get('.symbol-table').should('contain','0');
        cy.get('.symbol-table').should('contain','int');
        cy.get('.symbol-table').should('contain','red');
        cy.get('.symbol-table').should('contain','1');
        cy.get('.symbol-table').should('contain','int');
        cy.get('.symbol-table').should('contain','new');
        cy.get('.symbol-table').should('contain','0');
        cy.get('.symbol-table').should('contain','int');
        cy.get('.symbol-table').should('contain','_i');
        cy.get('.symbol-table').should('contain','0');

        cy.contains('Symbol table generated successfully');
    })


});*/
