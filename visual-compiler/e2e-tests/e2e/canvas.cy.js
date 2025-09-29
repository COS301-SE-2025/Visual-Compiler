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

    it('Add and connect Nodes', () => {

        cy.visit('http://localhost:5173/main-workspace');
        cy.wait(1000);

        const project_name = 'canvas_project';
        //delete project
        cy.get('.section-heading').should('contain', 'Start a new project');
        cy.get('.project-block').contains(project_name).get('.delete-button').click();
        cy.get('.delete-confirm-button').click();
        cy.wait(500);

        //select project
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
        cy.get('.code-input-container').should('contain', 'Source Code Input');
        cy.get('button').contains('Show Example').click();
        cy.get('.confirm-btn').contains('Confirm Code');
        cy.get('.confirm-btn').click();
        cy.contains('Source code saved');;

        //lexer node
        cy.get('button').should('contain', 'Lexer');
        cy.get('button').contains('Lexer').click();
        cy.get('.canvas-container').should('contain','Lexer');
        cy.get('#A-1\\/N-source-1').trigger('mousedown', {which: 1, force: true});
        cy.get('#A-1\\/N-lexer-2').trigger('mousemove', {force: true}).trigger('mouseup', {force: true});
        cy.wait(500);
        
        //parser node
        cy.get('button').should('contain', 'Parser');
        cy.get('button').contains('Parser').click();
        cy.get('.canvas-container').should('contain','Parser');
        cy.get('#A-2\\/N-lexer-2').trigger('mousedown', {which: 1, force: true});
        cy.get('#A-1\\/N-parser-3').trigger('mousemove', {force: true}).trigger('mouseup', {force: true});
        cy.wait(500);

        //analyser node
        cy.get('button').should('contain', 'Analyser');
        cy.get('button').contains('Analyser').click();
        cy.get('.canvas-container').should('contain','Analyser');
        cy.get('#A-2\\/N-parser-3').trigger('mousedown', {which: 1, force: true});
        cy.get('#A-1\\/N-analyser-4').trigger('mousemove', {force: true}).trigger('mouseup', {force: true});
        cy.wait(500);

        //translator node
        cy.get('button').should('contain', 'Translator');
        cy.get('button').contains('Translator').click();
        cy.get('.canvas-container').should('contain','Translator');
        cy.get('#A-2\\/N-analyser-4').trigger('mousedown', {which: 1, force: true});
        cy.get('#A-1\\/N-translator-5').trigger('mousemove', {force: true}).trigger('mouseup', {force: true});

        //optimiser node
        cy.get('button').should('contain', 'Optimiser');
        cy.get('button').contains('Optimiser').click();
        cy.get('.canvas-container').should('contain','Optimiser');

    })


})
