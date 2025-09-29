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

    it('No API connection', () => {

        cy.visit('http://localhost:5173/main-workspace');


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
		cy.get('button').contains("Next").click();
        cy.get('button').contains("Get Started").click();

        cy.wait(5000);

        //source code node
        cy.get('button').should('contain', 'Source Code');
        cy.get('button').contains('Source Code').click();
        cy.get('.canvas-container').should('contain','Source Code');
        cy.get('.canvas-container').contains('Source Code').dblclick();
        cy.wait(500);
        cy.get('.code-input-container').should('contain', 'Source Code Input');
        cy.contains('Failed to load projects. Please try again later.');
        cy.get('.code-input-header-row button').click();
        cy.get('.confirm-btn').contains('Confirm Code');
        cy.get('.confirm-btn').click();
        cy.contains('Save failed:');
  
    })


});
