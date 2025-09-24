describe('Lexer Test', ()=> {
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

    it('Lexer phase - regex', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        const project_name = `canvas_project`;
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
        cy.get('button').contains("Get Started").click();

        //source code node
        cy.get('button').should('contain', 'Source Code');
        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').should('contain','Source Code');
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.wait(500);
        cy.get('.code-input-container').should('contain', 'Source Code Input');
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').contains('Confirm Code');
        cy.get('.confirm-btn').click();
        cy.contains('Source code saved');

        //lexer node
        cy.get('button').should('contain', 'Lexer');
        cy.get('button').contains('Lexer').click();
        cy.get('.canvas-container').should('contain','Lexer');
        cy.get('#A-1\\/N-source-1').trigger('mousedown', { which: 1, force: true });
        cy.get('#A-1\\/N-lexer-2').trigger('mousemove', { force: true }).trigger('mouseup', { force: true });
        cy.get('.canvas-container').contains('Lexer').dblclick();
        cy.wait(500);

        cy.get('.phase-inspector').should('contain','LEXING');
        cy.get('.automaton-btn').contains('Regular Expression').click();
        cy.get('.default-toggle-btn').click();

        cy.get('.submit-button').contains('Submit').click();
        cy.get('.generate-button').contains('Generate Tokens').click();

        cy.get('.token-table').should('contain','Type');
        cy.get('.token-table').should('contain','Value');
        cy.get('.token-table').should('contain','KEYWORD');
        cy.get('.token-table').should('contain','int');
        cy.get('.token-table').should('contain','print');
        cy.get('.token-table').should('contain','IDENTIFIER');
        cy.get('.token-table').should('contain','blue');
        cy.get('.token-table').should('contain','new');
        cy.get('.token-table').should('contain','red');
        cy.get('.token-table').should('contain','_i');
        cy.get('.token-table').should('contain','ASSIGNMENT');
        cy.get('.token-table').should('contain','=');
        cy.get('.token-table').should('contain','INTEGER');
        cy.get('.token-table').should('contain','13');
        cy.get('.token-table').should('contain','1');
        cy.get('.token-table').should('contain','12');
        cy.get('.token-table').should('contain','DELIMITER');
        cy.get('.token-table').should('contain',';');
        cy.get('.token-table').should('contain','OPEN_BRACKET');
        cy.get('.token-table').should('contain','(');
        cy.get('.token-table').should('contain','CLOSE_BRACKET');
        cy.get('.token-table').should('contain',')');
        cy.get('.token-table').should('contain','OPERATOR');
        cy.get('.token-table').should('contain','+');
        cy.get('.token-table').should('contain','CONTROL');
        cy.get('.token-table').should('contain','range');

        cy.contains('Successfully tokenised');
  
    })

    it('Lexer phase - automata', () => {

        cy.visit('http://localhost:5173/main-workspace');

        cy.wait(1000);

        const project_name = `canvas_project`;
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
        cy.get('button').contains("Get Started").click();

        //source code node
        cy.get('button').should('contain', 'Source Code');
        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').should('contain','Source Code');
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.wait(500);
        cy.get('.code-input-container').should('contain', 'Source Code Input');
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').contains('Confirm Code');
        cy.get('.confirm-btn').click();
        cy.contains('Source code saved');

        //lexer node
        cy.get('button').should('contain', 'Lexer');
        cy.get('button').contains('Lexer').click();
        cy.get('.canvas-container').should('contain','Lexer');
        cy.get('#A-1\\/N-source-1').trigger('mousedown', { which: 1, force: true });
        cy.get('#A-1\\/N-lexer-2').trigger('mousemove', { force: true }).trigger('mouseup', { force: true });
        cy.get('.canvas-container').contains('Lexer').dblclick();
        cy.wait(500);

        cy.get('.phase-inspector').should('contain','LEXING');
        cy.get('.automaton-btn').contains('Automata').click();
        cy.get('.default-toggle-btn').click();

        cy.get('.action-btn').contains('Tokenisation').click();

        cy.get('.token-table').should('contain','Type');
        cy.get('.token-table').should('contain','Value');
        cy.get('.token-table').should('contain','KEYWORD');
        cy.get('.token-table').should('contain','int');
        cy.get('.token-table').should('contain','print');
        cy.get('.token-table').should('contain','IDENTIFIER');
        cy.get('.token-table').should('contain','blue');
        cy.get('.token-table').should('contain','new');
        cy.get('.token-table').should('contain','red');
        cy.get('.token-table').should('contain','_i');
        cy.get('.token-table').should('contain','ASSIGNMENT');
        cy.get('.token-table').should('contain','=');
        cy.get('.token-table').should('contain','INTEGER');
        cy.get('.token-table').should('contain','13');
        cy.get('.token-table').should('contain','1');
        cy.get('.token-table').should('contain','12');
        cy.get('.token-table').should('contain','DELIMITER');
        cy.get('.token-table').should('contain',';');
        cy.get('.token-table').should('contain','OPEN_BRACKET');
        cy.get('.token-table').should('contain','(');
        cy.get('.token-table').should('contain','CLOSE_BRACKET');
        cy.get('.token-table').should('contain',')');
        cy.get('.token-table').should('contain','OPERATOR');
        cy.get('.token-table').should('contain','+');
        cy.get('.token-table').should('contain','CONTROL');
        cy.get('.token-table').should('contain','range');

        cy.contains('Tokenization complete');
  
    })

});